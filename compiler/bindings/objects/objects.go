// Package objects provides object placement, scatter, and raycast for CyberBasic.
package objects

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"cyberbasic/compiler/vm"
)

// ObjectInstance holds model id and transform (position, scale, rotation).
type ObjectInstance struct {
	ModelID  string
	X, Y, Z  float32
	ScaleX   float32
	ScaleY   float32
	ScaleZ   float32
	RotAxisX float32
	RotAxisY float32
	RotAxisZ float32
	RotAngle float32
}

// ObjectExport is the JSON-serializable form for world save/load.
type ObjectExport struct {
	ModelID  string  `json:"modelId"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z"`
	ScaleX   float64 `json:"scaleX"`
	ScaleY   float64 `json:"scaleY"`
	ScaleZ   float64 `json:"scaleZ"`
	RotAxisX float64 `json:"rotAxisX"`
	RotAxisY float64 `json:"rotAxisY"`
	RotAxisZ float64 `json:"rotAxisZ"`
	RotAngle float64 `json:"rotAngle"`
}

var (
	objectInstances   = make(map[string]*ObjectInstance)
	objectSeq         int
	objectMu          sync.Mutex
)

func toFloat32(v interface{}) float32 {
	switch x := v.(type) {
	case int:
		return float32(x)
	case int32:
		return float32(x)
	case float64:
		return float32(x)
	case float32:
		return x
	default:
		return 0
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func toInt32(v interface{}) int32 {
	switch x := v.(type) {
	case int:
		return int32(x)
	case int32:
		return x
	case float64:
		return int32(x)
	default:
		return 0
	}
}

// RegisterObjects registers object placement and query bindings with the VM.
func RegisterObjects(v *vm.VM) {
	v.RegisterForeign("ObjectPlace", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("ObjectPlace requires (modelId, x, y, z, scale, rotation) or (modelId, x, y, z, scaleX, scaleY, scaleZ, rotAxisX, rotAxisY, rotAxisZ, rotAngle)")
		}
		modelID := toString(args[0])
		x, y, z := toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		sx, sy, sz := toFloat32(args[4]), toFloat32(args[4]), toFloat32(args[4])
		rax, ray, raz, angle := float32(0), float32(1), float32(0), float32(0)
		if len(args) >= 11 {
			sx, sy, sz = toFloat32(args[4]), toFloat32(args[5]), toFloat32(args[6])
			rax, ray, raz, angle = toFloat32(args[7]), toFloat32(args[8]), toFloat32(args[9]), toFloat32(args[10])
		} else if len(args) >= 6 {
			angle = toFloat32(args[5])
		}
		objectMu.Lock()
		objectSeq++
		id := fmt.Sprintf("obj_%d", objectSeq)
		objectInstances[id] = &ObjectInstance{
			ModelID:  modelID,
			X: x, Y: y, Z: z,
			ScaleX: sx, ScaleY: sy, ScaleZ: sz,
			RotAxisX: rax, RotAxisY: ray, RotAxisZ: raz, RotAngle: angle,
		}
		objectMu.Unlock()
		return id, nil
	})

	v.RegisterForeign("ObjectRemove", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ObjectRemove requires (objectId)")
		}
		id := toString(args[0])
		objectMu.Lock()
		_, ok := objectInstances[id]
		delete(objectInstances, id)
		objectMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id: %s", id)
		}
		return nil, nil
	})

	v.RegisterForeign("ObjectSetTransform", func(args []interface{}) (interface{}, error) {
		if len(args) < 12 {
			return nil, fmt.Errorf("ObjectSetTransform requires (objectId, x, y, z, scaleX, scaleY, scaleZ, rotAxisX, rotAxisY, rotAxisZ, rotAngle)")
		}
		id := toString(args[0])
		objectMu.Lock()
		o, ok := objectInstances[id]
		objectMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown object id: %s", id)
		}
		o.X, o.Y, o.Z = toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3])
		o.ScaleX, o.ScaleY, o.ScaleZ = toFloat32(args[4]), toFloat32(args[5]), toFloat32(args[6])
		o.RotAxisX, o.RotAxisY, o.RotAxisZ = toFloat32(args[7]), toFloat32(args[8]), toFloat32(args[9])
		o.RotAngle = toFloat32(args[10])
		return nil, nil
	})

	v.RegisterForeign("ObjectRandomScatter", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("ObjectRandomScatter requires (modelId, areaX, areaZ, count, minScale, maxScale)")
		}
		modelID := toString(args[0])
		areaX, areaZ := toFloat32(args[1]), toFloat32(args[2])
		count := int(toInt32(args[3]))
		minS, maxS := toFloat32(args[4]), toFloat32(args[5])
		if count <= 0 {
			count = 10
		}
		var ids []interface{}
		for i := 0; i < count; i++ {
			x := (float32(rand.Float64()) - 0.5) * 2 * areaX
			z := (float32(rand.Float64()) - 0.5) * 2 * areaZ
			scale := minS + float32(rand.Float64())*(maxS-minS)
			rot := float32(rand.Float64() * 2 * math.Pi)
			res, err := v.CallForeign("ObjectPlace", []interface{}{modelID, x, 0, z, scale, rot})
			if err != nil {
				continue
			}
			if id, ok := res.(string); ok {
				ids = append(ids, id)
			}
		}
		return ids, nil
	})

	v.RegisterForeign("ObjectPaint", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ObjectPaint requires (modelId, x, z, radius, density)")
		}
		modelID := toString(args[0])
		x, z := toFloat32(args[1]), toFloat32(args[2])
		radius, density := toFloat32(args[3]), toFloat32(args[4])
		n := int(density)
		if n <= 0 {
			n = 5
		}
		for i := 0; i < n; i++ {
			angle := float32(rand.Float64() * 2 * math.Pi)
			r := float32(rand.Float64()) * radius
			gx := x + r*float32(math.Cos(float64(angle)))
			gz := z + r*float32(math.Sin(float64(angle)))
			scale := 0.8 + float32(rand.Float64())*0.4
			rot := float32(rand.Float64() * 2 * math.Pi)
			_, _ = v.CallForeign("ObjectPlace", []interface{}{modelID, gx, 0, gz, scale, rot})
		}
		return nil, nil
	})

	v.RegisterForeign("ObjectErase", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ObjectErase requires (x, z, radius)")
		}
		x, z, radius := toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2])
		radiusSq := radius * radius
		objectMu.Lock()
		for id, o := range objectInstances {
			dx := o.X - x
			dz := o.Z - z
			if dx*dx+dz*dz <= radiusSq {
				delete(objectInstances, id)
			}
		}
		objectMu.Unlock()
		return nil, nil
	})

	v.RegisterForeign("ObjectGetAt", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ObjectGetAt requires (x, z)")
		}
		x, z := toFloat32(args[0]), toFloat32(args[1])
		objectMu.Lock()
		var best string
		bestDistSq := float32(math.MaxFloat32)
		for id, o := range objectInstances {
			dx := o.X - x
			dz := o.Z - z
			d2 := dx*dx + dz*dz
			if d2 < bestDistSq {
				bestDistSq = d2
				best = id
			}
		}
		objectMu.Unlock()
		return best, nil
	})

	v.RegisterForeign("ObjectRaycast", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("ObjectRaycast requires (ox, oy, oz, dx, dy, dz)")
		}
		ox := toFloat32(args[0])
		oy := toFloat32(args[1])
		oz := toFloat32(args[2])
		dx := toFloat32(args[3])
		dy := toFloat32(args[4])
		dz := toFloat32(args[5])
		// Simplified: return first object whose bounding sphere is hit
		objectMu.Lock()
		var hitID string
		var hitDist float32 = 1e9
		for id, o := range objectInstances {
			// Sphere at (o.X, o.Y, o.Z) with radius ~max(scale)
			r := o.ScaleX
			if o.ScaleY > r {
				r = o.ScaleY
			}
			if o.ScaleZ > r {
				r = o.ScaleZ
			}
			r *= 2
			// Ray-sphere intersection
			px := ox - o.X
			py := oy - o.Y
			pz := oz - o.Z
			a := dx*dx + dy*dy + dz*dz
			b := 2 * (px*dx + py*dy + pz*dz)
			c := px*px + py*py + pz*pz - r*r
			disc := b*b - 4*a*c
			if disc >= 0 {
				t := (-b - float32(math.Sqrt(float64(disc)))) / (2 * a)
				if t > 0 && t < hitDist {
					hitDist = t
					hitID = id
				}
			}
		}
		objectMu.Unlock()
		if hitID == "" {
			return []interface{}{0, "", 0.0, 0.0, 0.0}, nil
		}
		return []interface{}{1, hitID, ox + hitDist*dx, oy + hitDist*dy, oz + hitDist*dz}, nil
	})

	// DrawObject draws a single object by id.
	v.RegisterForeign("DrawObject", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawObject requires (objectId)")
		}
		id := toString(args[0])
		objectMu.Lock()
		o, ok := objectInstances[id]
		objectMu.Unlock()
		if !ok {
			return nil, nil
		}
		_, err := v.CallForeign("DrawModelEx", []interface{}{
			o.ModelID,
			o.X, o.Y, o.Z,
			o.RotAxisX, o.RotAxisY, o.RotAxisZ, o.RotAngle,
			o.ScaleX, o.ScaleY, o.ScaleZ,
		})
		return nil, err
	})

	// DrawAllObjects draws every placed object.
	v.RegisterForeign("DrawAllObjects", func(args []interface{}) (interface{}, error) {
		objectMu.Lock()
		list := make([]*ObjectInstance, 0, len(objectInstances))
		for _, o := range objectInstances {
			list = append(list, o)
		}
		objectMu.Unlock()
		for _, o := range list {
			_, _ = v.CallForeign("DrawModelEx", []interface{}{
				o.ModelID,
				o.X, o.Y, o.Z,
				o.RotAxisX, o.RotAxisY, o.RotAxisZ, o.RotAngle,
				o.ScaleX, o.ScaleY, o.ScaleZ,
			})
		}
		return nil, nil
	})
}

// ExportForSave returns a snapshot of all object instances for world save.
func ExportForSave() map[string]ObjectExport {
	objectMu.Lock()
	defer objectMu.Unlock()
	out := make(map[string]ObjectExport, len(objectInstances))
	for id, o := range objectInstances {
		out[id] = ObjectExport{
			ModelID:  o.ModelID,
			X:        float64(o.X), Y: float64(o.Y), Z: float64(o.Z),
			ScaleX:   float64(o.ScaleX), ScaleY: float64(o.ScaleY), ScaleZ: float64(o.ScaleZ),
			RotAxisX: float64(o.RotAxisX), RotAxisY: float64(o.RotAxisY), RotAxisZ: float64(o.RotAxisZ),
			RotAngle: float64(o.RotAngle),
		}
	}
	return out
}

// ImportFromLoad restores object instances from a world load (clears existing first).
func ImportFromLoad(data map[string]ObjectExport) {
	objectMu.Lock()
	defer objectMu.Unlock()
	objectInstances = make(map[string]*ObjectInstance)
	for id, e := range data {
		objectInstances[id] = &ObjectInstance{
			ModelID:  e.ModelID,
			X:        float32(e.X), Y: float32(e.Y), Z: float32(e.Z),
			ScaleX:   float32(e.ScaleX), ScaleY: float32(e.ScaleY), ScaleZ: float32(e.ScaleZ),
			RotAxisX: float32(e.RotAxisX), RotAxisY: float32(e.RotAxisY), RotAxisZ: float32(e.RotAxisZ),
			RotAngle: float32(e.RotAngle),
		}
	}
}
