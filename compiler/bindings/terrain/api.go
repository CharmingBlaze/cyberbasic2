package terrain

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
)

// TerrainState holds mesh, material, and heightmap refs for a terrain instance.
type TerrainState struct {
	HeightmapID    string
	MeshID         string
	MaterialID     string
	SizeX          float32
	SizeZ          float32
	HeightScale    float32
	LODLevel       int
	CollisionEnabled bool
	Friction      float32
	Bounce        float32
}

var (
	terrains   = make(map[string]*TerrainState)
	terrainSeq int
	terrainMu  sync.Mutex
	defaultMat string
)

// TerrainCreate runs the mesh generator and registers a terrain. Returns terrain id.
func TerrainCreate(v *vm.VM, heightmapID string, sizeX, sizeZ, heightScale float32) (string, error) {
	meshID, err := GenTerrainMesh(v, heightmapID, sizeX, sizeZ, heightScale, 0)
	if err != nil {
		return "", err
	}
	terrainMu.Lock()
	terrainSeq++
	id := fmt.Sprintf("terrain_%d", terrainSeq)
	terrains[id] = &TerrainState{
		HeightmapID: heightmapID,
		MeshID:      meshID,
		MaterialID:  "",
		SizeX:       sizeX,
		SizeZ:       sizeZ,
		HeightScale: heightScale,
		LODLevel:    0,
	}
	terrainMu.Unlock()
	return id, nil
}

// TerrainUpdate rebuilds the mesh from the current heightmap.
func TerrainUpdate(v *vm.VM, terrainID string) error {
	terrainMu.Lock()
	ts, ok := terrains[terrainID]
	terrainMu.Unlock()
	if !ok {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	meshID, err := GenTerrainMesh(v, ts.HeightmapID, ts.SizeX, ts.SizeZ, ts.HeightScale, ts.LODLevel)
	if err != nil {
		return err
	}
	terrainMu.Lock()
	ts.MeshID = meshID
	terrainMu.Unlock()
	return nil
}

// DrawTerrain draws the terrain at the given position (calls DrawMesh).
func DrawTerrain(v *vm.VM, terrainID string, posX, posY, posZ float32) error {
	terrainMu.Lock()
	ts, ok := terrains[terrainID]
	terrainMu.Unlock()
	if !ok {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	matID := ts.MaterialID
	if matID == "" {
		if defaultMat == "" {
			res, err := v.CallForeign("LoadMaterialDefault", nil)
			if err != nil || res == nil {
				matID = "default"
			} else {
				defaultMat, _ = res.(string)
				matID = defaultMat
			}
		} else {
			matID = defaultMat
		}
	}
	_, err := v.CallForeign("DrawMesh", []interface{}{
		ts.MeshID, matID,
		posX, posY, posZ,
		float32(1), float32(1), float32(1),
	})
	return err
}

// SetTerrainMaterial sets the material id used when drawing.
func SetTerrainMaterial(terrainID, materialID string) error {
	terrainMu.Lock()
	defer terrainMu.Unlock()
	ts, ok := terrains[terrainID]
	if !ok {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	ts.MaterialID = materialID
	return nil
}

// SetTerrainTexture is an alias that sets the terrain's material from a texture (user can create material first).
func SetTerrainTexture(terrainID, textureID string) error {
	return SetTerrainMaterial(terrainID, textureID)
}

// SetTerrainLOD sets LOD level; next TerrainUpdate will use it.
func SetTerrainLOD(terrainID string, lodLevel int) error {
	terrainMu.Lock()
	defer terrainMu.Unlock()
	ts, ok := terrains[terrainID]
	if !ok {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	ts.LODLevel = lodLevel
	return nil
}

// GetTerrainState returns internal state for a terrain (for edit/query).
func GetTerrainState(terrainID string) *TerrainState {
	terrainMu.Lock()
	ts := terrains[terrainID]
	terrainMu.Unlock()
	return ts
}
