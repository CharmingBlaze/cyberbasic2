package vegetation

import (
	"fmt"

	"cyberbasic/compiler/vm"
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

// RegisterVegetation registers tree and grass bindings with the VM.
func RegisterVegetation(v *vm.VM) {
	// Trees
	v.RegisterForeign("TreeTypeCreate", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("TreeTypeCreate requires (modelId, trunkTextureId, leafTextureId)")
		}
		return TreeTypeCreate(toString(args[0]), toString(args[1]), toString(args[2])), nil
	})
	v.RegisterForeign("TreeSystemCreate", func(args []interface{}) (interface{}, error) {
		return TreeSystemCreate(), nil
	})
	v.RegisterForeign("TreePlace", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("TreePlace requires (systemId, typeId, x, y, z, scale, rotation)")
		}
		return TreePlace(toString(args[0]), toString(args[1]), toFloat32(args[2]), toFloat32(args[3]), toFloat32(args[4]), toFloat32(args[5]), toFloat32(args[6]))
	})
	v.RegisterForeign("TreeRemove", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("TreeRemove requires (treeId)")
		}
		return nil, TreeRemove(toString(args[0]))
	})
	v.RegisterForeign("TreeSetPosition", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("TreeSetPosition requires (treeId, x, y, z)")
		}
		return nil, TreeSetPosition(toString(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
	})
	v.RegisterForeign("TreeSetScale", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("TreeSetScale requires (treeId, scale)")
		}
		return nil, TreeSetScale(toString(args[0]), toFloat32(args[1]))
	})
	v.RegisterForeign("TreeSetRotation", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("TreeSetRotation requires (treeId, rotation)")
		}
		return nil, TreeSetRotation(toString(args[0]), toFloat32(args[1]))
	})
	v.RegisterForeign("TreeSystemSetLOD", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("TreeSystemSetLOD requires (systemId, near, mid, far)")
		}
		TreeSystemSetLOD(toString(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
		return nil, nil
	})
	v.RegisterForeign("TreeSystemEnableInstancing", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("TreeSystemEnableInstancing requires (systemId, on)")
		}
		TreeSystemEnableInstancing(toString(args[0]), toFloat32(args[1]) != 0)
		return nil, nil
	})
	v.RegisterForeign("TreeGetAt", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("TreeGetAt requires (systemId, x, z)")
		}
		return TreeGetAt(toString(args[0]), toFloat32(args[1]), toFloat32(args[2])), nil
	})
	v.RegisterForeign("TreeEnableCollision", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("TreeEnableCollision requires (treeId, flag)")
		}
		TreeSetCollisionEnabled(toString(args[0]), toFloat32(args[1]) != 0)
		return nil, nil
	})
	v.RegisterForeign("TreeSetCollisionRadius", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("TreeSetCollisionRadius requires (treeId, radius)")
		}
		TreeSetCollisionRadius(toString(args[0]), toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("TreeSetWind", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("TreeSetWind requires (typeId, strength, speed)")
		}
		TreeTypeSetWind(toString(args[0]), toFloat32(args[1]), toFloat32(args[2]))
		return nil, nil
	})
	v.RegisterForeign("TreeApplyWind", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("TreeApplyWind requires (treeId, vx, vy, vz)")
		}
		// Stub: apply wind vector to tree (shader bending)
		return nil, nil
	})
	v.RegisterForeign("TreeRaycast", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("TreeRaycast requires (systemId, ox, oy, oz, dx, dy, dz)")
		}
		// Stub: ray vs tree capsules
		return 0, nil
	})

	v.RegisterForeign("DrawTrees", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawTrees requires (systemId)")
		}
		sysID := toString(args[0])
		ids := GetTreeSystemInstanceIds(sysID)
		if ids == nil {
			return nil, nil
		}
		for _, tid := range ids {
			t := getTreeInstance(tid)
			if t == nil {
				continue
			}
			tt := GetTreeType(t.TypeID)
			if tt == nil {
				continue
			}
			_, _ = v.CallForeign("DrawModelEx", []interface{}{
				tt.ModelID,
				t.X, t.Y, t.Z,
				float32(0), float32(1), float32(0), t.Rotation,
				t.Scale, t.Scale, t.Scale,
			})
		}
		return nil, nil
	})
	v.RegisterRenderType("drawtrees", vm.Render3D)

	// Grass
	v.RegisterForeign("GrassCreate", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GrassCreate requires (textureId, density, patchSize)")
		}
		return GrassCreate(toString(args[0]), toFloat32(args[1]), toFloat32(args[2])), nil
	})
	v.RegisterForeign("GrassSetWind", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GrassSetWind requires (grassId, speed, strength)")
		}
		GrassSetWind(toString(args[0]), toFloat32(args[1]), toFloat32(args[2]))
		return nil, nil
	})
	v.RegisterForeign("GrassSetHeight", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GrassSetHeight requires (grassId, height)")
		}
		GrassSetHeight(toString(args[0]), toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("GrassSetColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("GrassSetColor requires (grassId, r, g, b, a)")
		}
		GrassSetColor(toString(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]), toFloat32(args[4]))
		return nil, nil
	})
	v.RegisterForeign("GrassPaint", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("GrassPaint requires (grassId, x, z, radius, density)")
		}
		GrassPaint(toString(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]), toFloat32(args[4]))
		return nil, nil
	})
	v.RegisterForeign("GrassErase", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("GrassErase requires (grassId, x, z, radius)")
		}
		GrassErase(toString(args[0]), toFloat32(args[1]), toFloat32(args[2]), toFloat32(args[3]))
		return nil, nil
	})
	v.RegisterForeign("GrassSetDensity", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GrassSetDensity requires (grassId, density)")
		}
		GrassSetDensity(toString(args[0]), toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("GrassSetLOD", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GrassSetLOD requires (grassId, dist)")
		}
		GrassSetLOD(toString(args[0]), toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("GrassEnableInstancing", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GrassEnableInstancing requires (grassId, on)")
		}
		GrassEnableInstancing(toString(args[0]), toFloat32(args[1]) != 0)
		return nil, nil
	})
	v.RegisterForeign("GrassSetBendAmount", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GrassSetBendAmount requires (grassId, value)")
		}
		GrassSetBendAmount(toString(args[0]), toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("GrassSetInteraction", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GrassSetInteraction requires (grassId, flag)")
		}
		GrassSetInteraction(toString(args[0]), toFloat32(args[1]) != 0)
		return nil, nil
	})

	v.RegisterForeign("DrawGrass", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("DrawGrass requires (grassId)")
		}
		g := getGrass(toString(args[0]))
		if g == nil {
			return nil, nil
		}
		for _, inst := range g.Instances {
			_, _ = v.CallForeign("DrawBillboard", []interface{}{
				g.TextureID,
				inst.X, inst.Y, inst.Z,
				inst.Scale,
				int(255*g.ColorR), int(255*g.ColorG), int(255*g.ColorB), int(255*g.ColorA),
			})
		}
		return nil, nil
	})
	v.RegisterRenderType("drawgrass", vm.Render3D)
}
