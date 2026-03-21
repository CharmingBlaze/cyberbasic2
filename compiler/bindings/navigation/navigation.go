// Package navigation provides NavGrid, NavMesh, and NavAgent for pathfinding.
package navigation

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"cyberbasic/compiler/bindings/modfacade"
	"cyberbasic/compiler/bindings/terrain"
	"cyberbasic/compiler/vm"
)

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func toFloat64(v interface{}) float64 {
	switch x := v.(type) {
	case int:
		return float64(x)
	case float64:
		return x
	case float32:
		return float64(x)
	default:
		return 0
	}
}

func toInt(v interface{}) int {
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(math.Trunc(x))
	default:
		return 0
	}
}

type navGrid struct {
	width    int
	height   int
	walkable [][]bool
	cost     [][]float64
}

	var (
		grids   = make(map[string]*navGrid)
		gridSeq int
		gridsMu sync.RWMutex

		navMeshes   = make(map[string]*waypointGraph)
		navMeshSeq  int
		navMeshesMu sync.RWMutex

		navAgents   = make(map[string]*navAgent)
		navAgentSeq int
		navAgentsMu sync.RWMutex
	)

	type waypointGraph struct {
		verts     []struct{ x, y, z float64 }
		edges     map[int][]int
		obstacles []struct{ minX, minY, minZ, maxX, maxY, maxZ float64 }
	}

	type navAgent struct {
		meshId    string
		gridId    string
		x, y, z   float64
		destX, destY, destZ float64
		speed     float64
		radius    float64
		path      []struct{ x, y, z float64 }
		pathIndex int
	}

// RegisterNavigation registers NavGrid, NavMesh, NavAgent commands.
func RegisterNavigation(v *vm.VM) {
	v.RegisterForeign("NavGridCreate", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("NavGridCreate requires (width, height)")
		}
		w := toInt(args[0])
		h := toInt(args[1])
		if w <= 0 || h <= 0 {
			return nil, fmt.Errorf("NavGridCreate: width and height must be positive")
		}
		walkable := make([][]bool, w)
		cost := make([][]float64, w)
		for x := 0; x < w; x++ {
			walkable[x] = make([]bool, h)
			cost[x] = make([]float64, h)
			for y := 0; y < h; y++ {
				walkable[x][y] = true
				cost[x][y] = 1
			}
		}
		gridsMu.Lock()
		gridSeq++
		id := fmt.Sprintf("navgrid_%d", gridSeq)
		grids[id] = &navGrid{width: w, height: h, walkable: walkable, cost: cost}
		gridsMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("NavGridSetWalkable", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("NavGridSetWalkable requires (gridId, x, y, flag)")
		}
		gridId := toString(args[0])
		x, y := toInt(args[1]), toInt(args[2])
		flag := toFloat64(args[3]) != 0
		gridsMu.RLock()
		g := grids[gridId]
		gridsMu.RUnlock()
		if g == nil || x < 0 || x >= g.width || y < 0 || y >= g.height {
			return nil, nil
		}
		g.walkable[x][y] = flag
		return nil, nil
	})
	v.RegisterForeign("NavGridSetCost", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("NavGridSetCost requires (gridId, x, y, cost)")
		}
		gridId := toString(args[0])
		x, y := toInt(args[1]), toInt(args[2])
		c := toFloat64(args[3])
		if c < 0 {
			c = 0
		}
		gridsMu.RLock()
		g := grids[gridId]
		gridsMu.RUnlock()
		if g == nil || x < 0 || x >= g.width || y < 0 || y >= g.height {
			return nil, nil
		}
		g.cost[x][y] = c
		return nil, nil
	})
	v.RegisterForeign("NavGridFindPath", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("NavGridFindPath requires (gridId, startX, startY, endX, endY)")
		}
		gridId := toString(args[0])
		sx, sy := toInt(args[1]), toInt(args[2])
		ex, ey := toInt(args[3]), toInt(args[4])
		gridsMu.RLock()
		g := grids[gridId]
		gridsMu.RUnlock()
		if g == nil {
			return []interface{}{}, nil
		}
		path := navGridAStar(g, sx, sy, ex, ey)
		result := make([]interface{}, 0, len(path)*2)
		for _, p := range path {
			result = append(result, float64(p.x), float64(p.y))
		}
		return result, nil
	})
	v.RegisterForeign("NavMeshLoadFromFile", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("NavMeshLoadFromFile requires (path)")
		}
		path := toString(args[0])
		g, err := loadWaypointGraph(path)
		if err != nil {
			return nil, err
		}
		navMeshesMu.Lock()
		navMeshSeq++
		meshId := fmt.Sprintf("navmesh_%d", navMeshSeq)
		navMeshes[meshId] = g
		navMeshesMu.Unlock()
		return meshId, nil
	})
	v.RegisterForeign("NavMeshCreateFromTerrain", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("NavMeshCreateFromTerrain requires (terrainId)")
		}
		terrainID := toString(args[0])
		ts := terrain.GetTerrainState(terrainID)
		if ts == nil {
			return nil, fmt.Errorf("unknown terrain id: %s", terrainID)
		}
		sizeX := float64(ts.SizeX)
		sizeZ := float64(ts.SizeZ)
		if sizeX <= 0 {
			sizeX = 100
		}
		if sizeZ <= 0 {
			sizeZ = 100
		}
		gridRes := 24
		if len(args) >= 2 {
			if r := toInt(args[1]); r > 0 && r <= 128 {
				gridRes = r
			}
		}
		maxStep := 2.0
		if len(args) >= 3 {
			if s := toFloat64(args[2]); s > 0 {
				maxStep = s
			}
		}
		g := &waypointGraph{edges: make(map[int][]int)}
		stepX := sizeX / float64(gridRes)
		stepZ := sizeZ / float64(gridRes)
		ox := -sizeX / 2
		oz := -sizeZ / 2
		for iz := 0; iz <= gridRes; iz++ {
			for ix := 0; ix <= gridRes; ix++ {
				x := ox + float64(ix)*stepX
				z := oz + float64(iz)*stepZ
				y, err := terrain.TerrainGetHeight(terrainID, x, z)
				if err != nil {
					y = 0
				}
				g.verts = append(g.verts, struct{ x, y, z float64 }{x, y, z})
			}
		}
		// Build edges: 8-neighbor, connect if height delta < maxStep
		for iz := 0; iz <= gridRes; iz++ {
			for ix := 0; ix <= gridRes; ix++ {
				i := iz*(gridRes+1) + ix
				for diz := -1; diz <= 1; diz++ {
					for dix := -1; dix <= 1; dix++ {
						if dix == 0 && diz == 0 {
							continue
						}
						jx, jz := ix+dix, iz+diz
						if jx < 0 || jx > gridRes || jz < 0 || jz > gridRes {
							continue
						}
						j := jz*(gridRes+1) + jx
						vi, vj := g.verts[i], g.verts[j]
						dy := vi.y - vj.y
						if dy < 0 {
							dy = -dy
						}
						if dy <= maxStep {
							g.edges[i] = append(g.edges[i], j)
						}
					}
				}
			}
		}
		navMeshesMu.Lock()
		navMeshSeq++
		meshId := fmt.Sprintf("navmesh_%d", navMeshSeq)
		navMeshes[meshId] = g
		navMeshesMu.Unlock()
		return meshId, nil
	})
	v.RegisterForeign("NavMeshAddObstacle", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, nil
		}
		meshId := toString(args[0])
		ob := struct{ minX, minY, minZ, maxX, maxY, maxZ float64 }{
			toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3]),
			toFloat64(args[4]), toFloat64(args[5]), toFloat64(args[6]),
		}
		navMeshesMu.Lock()
		if g := navMeshes[meshId]; g != nil {
			g.obstacles = append(g.obstacles, ob)
		}
		navMeshesMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("NavMeshRemoveObstacle", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		meshId := toString(args[0])
		idx := toInt(args[1])
		navMeshesMu.Lock()
		if g := navMeshes[meshId]; g != nil && idx >= 0 && idx < len(g.obstacles) {
			g.obstacles = append(g.obstacles[:idx], g.obstacles[idx+1:]...)
		}
		navMeshesMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("NavMeshFindPathRaw", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("NavMeshFindPathRaw requires (meshId, ox, oy, oz, dx, dy, dz)")
		}
		meshId := toString(args[0])
		ox, oy, oz := toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])
		dx, dy, dz := toFloat64(args[4]), toFloat64(args[5]), toFloat64(args[6])
		navMeshesMu.RLock()
		g := navMeshes[meshId]
		navMeshesMu.RUnlock()
		if g == nil || len(g.verts) == 0 {
			return []interface{}{}, nil
		}
		path := navMeshAStar(g, ox, oy, oz, dx, dy, dz)
		result := make([]interface{}, 0, len(path)*3)
		for _, v := range path {
			result = append(result, v.x, v.y, v.z)
		}
		return result, nil
	})
	v.RegisterForeign("NavAgentCreate", func(args []interface{}) (interface{}, error) {
		meshId, gridId := "", ""
		if len(args) >= 1 {
			meshId = toString(args[0])
		}
		if len(args) >= 2 {
			gridId = toString(args[1])
		}
		navAgentsMu.Lock()
		navAgentSeq++
		id := fmt.Sprintf("navagent_%d", navAgentSeq)
		navAgents[id] = &navAgent{meshId: meshId, gridId: gridId, speed: 1, radius: 0.5}
		navAgentsMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("NavAgentSetSpeed", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		id := toString(args[0])
		navAgentsMu.Lock()
		if a := navAgents[id]; a != nil {
			a.speed = toFloat64(args[1])
		}
		navAgentsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("NavAgentSetRadius", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		id := toString(args[0])
		navAgentsMu.Lock()
		if a := navAgents[id]; a != nil {
			a.radius = toFloat64(args[1])
		}
		navAgentsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("NavAgentSetDestination", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, nil
		}
		id := toString(args[0])
		dx, dy, dz := toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])
		navAgentsMu.Lock()
		a := navAgents[id]
		navAgentsMu.Unlock()
		if a == nil {
			return nil, nil
		}
		a.destX, a.destY, a.destZ = dx, dy, dz
		var pathResult interface{}
		var err error
		if a.meshId != "" {
			pathResult, err = v.CallForeign("NavMeshFindPathRaw", []interface{}{
				a.meshId, a.x, a.y, a.z, dx, dy, dz,
			})
		} else if a.gridId != "" {
			sx, sy := int(a.x), int(a.y)
			ex, ey := int(dx), int(dy)
			pathResult, err = v.CallForeign("NavGridFindPath", []interface{}{
				a.gridId, float64(sx), float64(sy), float64(ex), float64(ey),
			})
		}
		if err != nil {
			return nil, err
		}
		a.path = nil
		a.pathIndex = 0
		if arr, ok := pathResult.([]interface{}); ok && len(arr) >= 3 {
			for i := 0; i+2 < len(arr); i += 3 {
				a.path = append(a.path, struct{ x, y, z float64 }{
					toFloat64(arr[i]), toFloat64(arr[i+1]), toFloat64(arr[i+2]),
				})
			}
		} else if arr, ok := pathResult.([]interface{}); ok && len(arr) >= 2 {
			for i := 0; i+1 < len(arr); i += 2 {
				a.path = append(a.path, struct{ x, y, z float64 }{
					toFloat64(arr[i]), 0, toFloat64(arr[i+1]),
				})
			}
		}
		return nil, nil
	})
	v.RegisterForeign("NavAgentGetNextWaypoint", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return []interface{}{0.0, 0.0, 0.0}, nil
		}
		id := toString(args[0])
		navAgentsMu.RLock()
		a := navAgents[id]
		navAgentsMu.RUnlock()
		if a == nil || a.pathIndex >= len(a.path) {
			return []interface{}{0.0, 0.0, 0.0}, nil
		}
		w := a.path[a.pathIndex]
		return []interface{}{w.x, w.y, w.z}, nil
	})
	v.RegisterForeign("NavAgentUpdate", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, nil
		}
		id := toString(args[0])
		dt := toFloat64(args[1])
		navAgentsMu.Lock()
		a := navAgents[id]
		navAgentsMu.Unlock()
		if a == nil || a.pathIndex >= len(a.path) {
			return nil, nil
		}
		w := a.path[a.pathIndex]
		dx, dy, dz := w.x-a.x, w.y-a.y, w.z-a.z
		dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
		if dist < 1e-6 {
			a.pathIndex++
			return nil, nil
		}
		move := a.speed * dt
		if move >= dist {
			a.x, a.y, a.z = w.x, w.y, w.z
			a.pathIndex++
		} else {
			a.x += dx * (move / dist)
			a.y += dy * (move / dist)
			a.z += dz * (move / dist)
		}
		return nil, nil
	})
	v.RegisterForeign("NavAgentSetPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, nil
		}
		id := toString(args[0])
		navAgentsMu.Lock()
		if a := navAgents[id]; a != nil {
			a.x, a.y, a.z = toFloat64(args[1]), toFloat64(args[2]), toFloat64(args[3])
		}
		navAgentsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("NavAgentGetPositionX", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		navAgentsMu.RLock()
		a := navAgents[toString(args[0])]
		navAgentsMu.RUnlock()
		if a == nil {
			return 0.0, nil
		}
		return a.x, nil
	})
	v.RegisterForeign("NavAgentGetPositionY", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		navAgentsMu.RLock()
		a := navAgents[toString(args[0])]
		navAgentsMu.RUnlock()
		if a == nil {
			return 0.0, nil
		}
		return a.y, nil
	})
	v.RegisterForeign("NavAgentGetPositionZ", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return 0.0, nil
		}
		navAgentsMu.RLock()
		a := navAgents[toString(args[0])]
		navAgentsMu.RUnlock()
		if a == nil {
			return 0.0, nil
		}
		return a.z, nil
	})

	v.SetGlobal("navigation", modfacade.New(v, MethodToForeign))
}

type gridCell struct {
	x, y int
}

func (c gridCell) key() int64 {
	return int64(c.x)<<32 | int64(c.y)
}

type astarNode struct {
	cell gridCell
	g    float64
	f    float64
}

func navGridAStar(g *navGrid, sx, sy, ex, ey int) []gridCell {
	if sx < 0 || sx >= g.width || sy < 0 || sy >= g.height ||
		ex < 0 || ex >= g.width || ey < 0 || ey >= g.height {
		return nil
	}
	if !g.walkable[sx][sy] || !g.walkable[ex][ey] {
		return nil
	}
	if sx == ex && sy == ey {
		return []gridCell{{sx, sy}}
	}

	open := &nodeHeap{}
	heap.Init(open)
	heap.Push(open, astarNode{gridCell{sx, sy}, 0, manhattan(sx, sy, ex, ey)})
	cameFrom := make(map[int64]gridCell)
	gScore := make(map[int64]float64)
	gScore[gridCell{sx, sy}.key()] = 0

	dirs := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}, {1, 1}, {1, -1}, {-1, -1}, {-1, 1}}

	for open.Len() > 0 {
		cur := heap.Pop(open).(astarNode)
		cx, cy := cur.cell.x, cur.cell.y
		if cx == ex && cy == ey {
			path := []gridCell{{ex, ey}}
			k := cur.cell.key()
			for {
				prev, ok := cameFrom[k]
				if !ok {
					break
				}
				path = append([]gridCell{prev}, path...)
				k = prev.key()
				if prev.x == sx && prev.y == sy {
					break
				}
			}
			return path
		}

		for _, d := range dirs {
			nx, ny := cx+d[0], cy+d[1]
			if nx < 0 || nx >= g.width || ny < 0 || ny >= g.height || !g.walkable[nx][ny] {
				continue
			}
			nc := gridCell{nx, ny}
			stepCost := g.cost[nx][ny]
			if d[0] != 0 && d[1] != 0 {
				stepCost *= 1.414
			}
			tentG := gScore[cur.cell.key()] + stepCost
			nk := nc.key()
			if prev, ok := gScore[nk]; ok && tentG >= prev {
				continue
			}
			cameFrom[nk] = cur.cell
			gScore[nk] = tentG
			heap.Push(open, astarNode{nc, tentG, tentG + manhattan(nx, ny, ex, ey)})
		}
	}
	return nil
}

func loadWaypointGraph(path string) (*waypointGraph, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	g := &waypointGraph{edges: make(map[int][]int)}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 3 {
			x, _ := strconv.ParseFloat(parts[0], 64)
			y, _ := strconv.ParseFloat(parts[1], 64)
			z, _ := strconv.ParseFloat(parts[2], 64)
			g.verts = append(g.verts, struct{ x, y, z float64 }{x, y, z})
		} else if len(parts) == 2 {
			i, _ := strconv.Atoi(parts[0])
			j, _ := strconv.Atoi(parts[1])
			if i >= 0 && i < len(g.verts) && j >= 0 && j < len(g.verts) {
				g.edges[i] = append(g.edges[i], j)
				g.edges[j] = append(g.edges[j], i)
			}
		}
	}
	return g, sc.Err()
}

func dist3(ax, ay, az, bx, by, bz float64) float64 {
	dx, dy, dz := bx-ax, by-ay, bz-az
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func navMeshAStar(g *waypointGraph, ox, oy, oz, dx, dy, dz float64) []struct{ x, y, z float64 } {
	if len(g.verts) == 0 {
		return nil
	}
	si := nearestVert(g, ox, oy, oz)
	ei := nearestVert(g, dx, dy, dz)
	if si < 0 || ei < 0 || si == ei {
		if si >= 0 {
			v := g.verts[si]
			return []struct{ x, y, z float64 }{{v.x, v.y, v.z}}
		}
		return nil
	}
	open := &meshNodeHeap{}
	heap.Init(open)
	heap.Push(open, meshNode{si, 0, dist3(ox, oy, oz, g.verts[ei].x, g.verts[ei].y, g.verts[ei].z)})
	cameFrom := make(map[int]int)
	gScore := make(map[int]float64)
	gScore[si] = 0
	dest := g.verts[ei]
	for open.Len() > 0 {
		cur := heap.Pop(open).(meshNode)
		if cur.i == ei {
			path := []struct{ x, y, z float64 }{dest}
			for cur.i != si {
				cur.i = cameFrom[cur.i]
				v := g.verts[cur.i]
				path = append([]struct{ x, y, z float64 }{{v.x, v.y, v.z}}, path...)
			}
			return path
		}
		curV := g.verts[cur.i]
		for _, ni := range g.edges[cur.i] {
			nv := g.verts[ni]
			tentG := gScore[cur.i] + dist3(curV.x, curV.y, curV.z, nv.x, nv.y, nv.z)
			if prev, ok := gScore[ni]; ok && tentG >= prev {
				continue
			}
			cameFrom[ni] = cur.i
			gScore[ni] = tentG
			heap.Push(open, meshNode{ni, tentG, tentG + dist3(nv.x, nv.y, nv.z, dest.x, dest.y, dest.z)})
		}
	}
	return nil
}

type meshNode struct {
	i int
	g float64
	f float64
}

func nearestVert(g *waypointGraph, x, y, z float64) int {
	best := -1
	bestD := math.MaxFloat64
	for i, v := range g.verts {
		d := dist3(x, y, z, v.x, v.y, v.z)
		if d < bestD {
			bestD = d
			best = i
		}
	}
	return best
}

type meshNodeHeap []meshNode

func (h meshNodeHeap) Len() int           { return len(h) }
func (h meshNodeHeap) Less(i, j int) bool { return h[i].f < h[j].f }
func (h meshNodeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *meshNodeHeap) Push(x interface{}) { *h = append(*h, x.(meshNode)) }
func (h *meshNodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func manhattan(ax, ay, bx, by int) float64 {
	dx := ax - bx
	dy := ay - by
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return float64(dx + dy)
}

type nodeHeap []astarNode

func (h nodeHeap) Len() int            { return len(h) }
func (h nodeHeap) Less(i, j int) bool  { return h[i].f < h[j].f }
func (h nodeHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *nodeHeap) Push(x interface{}) { *h = append(*h, x.(astarNode)) }
func (h *nodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
