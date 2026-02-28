// Package raylib: hybrid update/draw — render command registry, ClearRenderQueues, FlushRenderQueues.
package raylib

import (
	"sort"
	"strings"

	"cyberbasic/compiler/bindings/game"
	"cyberbasic/compiler/vm"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerHybrid(v *vm.VM) {
	registerRenderTypes(v)
	v.RegisterForeign("ClearRenderQueues", func(args []interface{}) (interface{}, error) {
		v.ClearRenderQueues()
		return nil, nil
	})
	v.RegisterForeign("FlushRenderQueues", func(args []interface{}) (interface{}, error) {
		return flushRenderQueues(v)
	})
	v.RegisterForeign("SpriteBatchBegin", func(args []interface{}) (interface{}, error) {
		spriteBatchActive = true
		spriteBatchList = spriteBatchList[:0]
		return nil, nil
	})
	v.RegisterForeign("SpriteBatchEnd", func(args []interface{}) (interface{}, error) {
		// Flush is done in the render loop when this command is executed
		return nil, nil
	})
}

type sortable2D struct {
	order  int
	zIndex int
	item   vm.RenderQueueItem
}

var (
	spriteBatchActive bool
	spriteBatchList   []vm.RenderQueueItem
)

func getTextureIDForBatch(name string, args []interface{}) string {
	n := strings.ToLower(name)
	if len(args) < 1 {
		return ""
	}
	switch n {
	case "spritedraw":
		id := toString(args[0])
		spriteMu.Lock()
		s := sprites[id]
		texID := ""
		if s != nil {
			texID = s.TextureId
		}
		spriteMu.Unlock()
		return texID
	case "drawtexture", "drawtextureex", "drawtexturerec", "drawtexturepro":
		return toString(args[0])
	default:
		return ""
	}
}

func isBatchableDraw(name string) bool {
	n := strings.ToLower(name)
	return n == "spritedraw" || n == "drawtexture" || n == "drawtextureex" || n == "drawtexturerec" || n == "drawtexturepro"
}

func flushSpriteBatch(v *vm.VM) {
	if len(spriteBatchList) == 0 {
		spriteBatchActive = false
		return
	}
	// Sort by texture ID for minimal state changes
	sort.Slice(spriteBatchList, func(i, j int) bool {
		ti := getTextureIDForBatch(spriteBatchList[i].Name, spriteBatchList[i].Args)
		tj := getTextureIDForBatch(spriteBatchList[j].Name, spriteBatchList[j].Args)
		return ti < tj
	})
	for _, item := range spriteBatchList {
		_, _ = v.CallForeign(item.Name, item.Args)
	}
	spriteBatchList = nil
	spriteBatchActive = false
}

func resolveLayerAndZ(name string, args []interface{}) (layerID string, zIndex int) {
	n := strings.ToLower(name)
	if len(args) < 1 {
		return "", 0
	}
	id := toString(args[0])
	switch n {
	case "spritedraw":
		return GetSpriteLayerAndZ(id)
	case "drawtilemap":
		return game.GetTilemapLayerAndZ(id)
	case "drawparticles":
		return game.GetParticleSystemLayerAndZ(id)
	case "drawparticleemitter":
		return GetParticleEmitter2DLayerAndZ(id)
	default:
		return "", 0
	}
}

func flushRenderQueues(v *vm.VM) (interface{}, error) {
	q2D, q3D, qGUI := v.GetRenderQueues()
	rl.BeginDrawing()
	rl.ClearBackground(rl.NewColor(25, 25, 35, 255))
	// 2D: sort by layer order then z-index, then draw per layer with parallax/scroll
	baseCam := getCurrentCamera2D()
	sorted := make([]sortable2D, 0, len(q2D))
	for _, item := range q2D {
		layerID, zIndex := resolveLayerAndZ(item.Name, item.Args)
		order := GetLayerOrder(layerID)
		sorted = append(sorted, sortable2D{order: order, zIndex: zIndex, item: item})
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].order != sorted[j].order {
			return sorted[i].order < sorted[j].order
		}
		return sorted[i].zIndex < sorted[j].zIndex
	})
	// Draw 2D by layer: for each distinct (order,layerID), if visible, apply camera and draw items
	var prevOrder int = -1
	var prevLayerID string = "\x00"
	var in2D bool
	for i := 0; i < len(sorted); i++ {
		layerID, _ := resolveLayerAndZ(sorted[i].item.Name, sorted[i].item.Args)
		order := sorted[i].order
		if order != prevOrder || layerID != prevLayerID {
			if in2D {
				rl.EndMode2D()
				in2D = false
			}
			prevOrder = order
			prevLayerID = layerID
			if GetLayerVisible(layerID) {
				px, py := GetLayerParallax(layerID)
				sx, sy := GetLayerScroll(layerID)
				cam := baseCam
				cam.Target.X = baseCam.Target.X*px + sx
				cam.Target.Y = baseCam.Target.Y*py + sy
				cam.Offset.X = baseCam.Offset.X + sx
				cam.Offset.Y = baseCam.Offset.Y + sy
				rl.BeginMode2D(cam)
				in2D = true
			}
		}
		if in2D {
			nm := strings.ToLower(sorted[i].item.Name)
			if nm == "spritebatchend" {
				flushSpriteBatch(v)
			} else if nm == "spritebatchbegin" {
				spriteBatchActive = true
				spriteBatchList = spriteBatchList[:0]
			} else if spriteBatchActive && isBatchableDraw(nm) {
				spriteBatchList = append(spriteBatchList, sorted[i].item)
			} else if nm == "spritedraw" && len(sorted[i].item.Args) >= 1 && !SpriteInCullingRect(toString(sorted[i].item.Args[0])) {
				// Skip 2D culling: sprite outside view
			} else {
				_, _ = v.CallForeign(sorted[i].item.Name, sorted[i].item.Args)
			}
		}
	}
	if in2D {
		rl.EndMode2D()
	}
	rl.BeginMode3D(camera3D)
	for _, item := range q3D {
		_, _ = v.CallForeign(item.Name, item.Args)
	}
	rl.EndMode3D()
	// GUI (raygui) draws in 2D context
	rl.BeginMode2D(baseCam)
	for _, item := range qGUI {
		_, _ = v.CallForeign(item.Name, item.Args)
	}
	rl.EndMode2D()
	rl.EndDrawing()
	return nil, nil
}

func registerRenderTypes(v *vm.VM) {
	reg := func(name string, typ vm.RenderType) {
		v.RegisterRenderType(strings.ToLower(name), typ)
	}
	// 2D: shapes, text, textures, clear, mode
	reg("DrawRectangle", vm.Render2D)
	reg("rect", vm.Render2D)
	reg("DrawCircle", vm.Render2D)
	reg("circle", vm.Render2D)
	reg("DrawLine", vm.Render2D)
	reg("DrawLineV", vm.Render2D)
	reg("DrawCircleLines", vm.Render2D)
	reg("DrawRectangleLines", vm.Render2D)
	reg("DrawTriangle", vm.Render2D)
	reg("DrawTriangleLines", vm.Render2D)
	reg("DrawPixel", vm.Render2D)
	reg("DrawPoly", vm.Render2D)
	reg("DrawEllipse", vm.Render2D)
	reg("DrawRing", vm.Render2D)
	reg("DrawRectangleRounded", vm.Render2D)
	reg("DrawFPS", vm.Render2D)
	reg("DrawLineEx", vm.Render2D)
	reg("DrawPixelV", vm.Render2D)
	reg("DrawCircleSector", vm.Render2D)
	reg("DrawCircleGradient", vm.Render2D)
	reg("DrawCircleV", vm.Render2D)
	reg("DrawEllipseLines", vm.Render2D)
	reg("DrawRingLines", vm.Render2D)
	reg("DrawRectangleV", vm.Render2D)
	reg("DrawRectangleRec", vm.Render2D)
	reg("DrawRectanglePro", vm.Render2D)
	reg("DrawRectangleLinesEx", vm.Render2D)
	reg("DrawRectangleRoundedLines", vm.Render2D)
	reg("DrawPolyLines", vm.Render2D)
	reg("DrawText", vm.Render2D)
	reg("DrawTextSimple", vm.Render2D)
	reg("DrawTextEx", vm.Render2D)
	reg("DrawTextPro", vm.Render2D)
	reg("DrawSprite", vm.Render2D)
	reg("DrawTexture", vm.Render2D)
	reg("sprite", vm.Render2D)
	reg("DrawTextureEx", vm.Render2D)
	reg("DrawTextureRec", vm.Render2D)
	reg("DrawTexturePro", vm.Render2D)
	reg("DrawTextureFlipH", vm.Render2D)
	reg("DrawTextureFlipV", vm.Render2D)
	reg("DrawTextureV", vm.Render2D)
	reg("DrawTextureNPatch", vm.Render2D)
	reg("SpriteDraw", vm.Render2D)
	reg("DrawTextExFont", vm.Render2D)
	reg("DrawTextCodepoint", vm.Render2D)
	reg("DrawTextCodepoints", vm.Render2D)
	reg("DrawSpriteAnimation", vm.Render2D)
	reg("DrawEntity", vm.Render2D)
	reg("DrawView", vm.Render2D)
	reg("ClearBackground", vm.Render2D)
	reg("Background", vm.Render2D)
	reg("BeginMode2D", vm.Render2D)
	reg("EndMode2D", vm.Render2D)
	reg("DrawTilemap", vm.Render2D)
	reg("DrawBackground", vm.Render2D)
	reg("SpriteBatchBegin", vm.Render2D)
	reg("SpriteBatchEnd", vm.Render2D)
	reg("DrawParticleEmitter", vm.Render2D)
	// 3D
	reg("BeginMode3D", vm.Render3D)
	reg("EndMode3D", vm.Render3D)
	reg("DrawGrid", vm.Render3D)
	reg("DrawModel", vm.Render3D)
	reg("DrawModelWithState", vm.Render3D)
	reg("DrawModelSimple", vm.Render3D)
	reg("DrawCube", vm.Render3D)
	reg("cube", vm.Render3D)
	reg("DrawCubeWires", vm.Render3D)
	reg("DrawSphere", vm.Render3D)
	reg("DrawSphereWires", vm.Render3D)
	reg("DrawPlane", vm.Render3D)
	reg("DrawLine3D", vm.Render3D)
	reg("DrawPoint3D", vm.Render3D)
	reg("DrawCircle3D", vm.Render3D)
	reg("DrawCubeV", vm.Render3D)
	reg("DrawCylinder", vm.Render3D)
	reg("DrawCylinderWires", vm.Render3D)
	reg("DrawRay", vm.Render3D)
	reg("DrawTriangle3D", vm.Render3D)
	reg("DrawTriangleStrip3D", vm.Render3D)
	reg("DrawCubeWiresV", vm.Render3D)
	reg("DrawSphereEx", vm.Render3D)
	reg("DrawCylinderEx", vm.Render3D)
	reg("DrawCylinderWiresEx", vm.Render3D)
	reg("DrawCapsule", vm.Render3D)
	reg("DrawCapsuleWires", vm.Render3D)
	reg("DrawModelEx", vm.Render3D)
	reg("DrawModelWires", vm.Render3D)
	reg("DrawBoundingBox", vm.Render3D)
	reg("DrawModelWiresEx", vm.Render3D)
	reg("DrawModelPoints", vm.Render3D)
	reg("DrawModelPointsEx", vm.Render3D)
	reg("DrawBillboard", vm.Render3D)
	reg("DrawBillboardRec", vm.Render3D)
	reg("DrawBillboardPro", vm.Render3D)
	reg("DrawText3D", vm.Render3D)
	reg("DrawMesh", vm.Render3D)
	reg("DrawMeshMatrix", vm.Render3D)
	reg("DrawMeshInstanced", vm.Render3D)
	reg("DrawTerrain", vm.Render3D)
	reg("DrawWater", vm.Render3D)
	reg("DrawTrees", vm.Render3D)
	reg("DrawGrass", vm.Render3D)
	reg("DrawObject", vm.Render3D)
	reg("DrawLevelObject", vm.Render3D)
	// GUI (raygui)
	reg("GuiLabel", vm.RenderGUI)
	reg("GuiButton", vm.RenderGUI)
	reg("button", vm.RenderGUI)
	reg("GuiCheckBox", vm.RenderGUI)
	reg("GuiCheckbox", vm.RenderGUI)
	reg("GuiSlider", vm.RenderGUI)
	reg("GuiProgressBar", vm.RenderGUI)
	reg("GuiTextbox", vm.RenderGUI)
	reg("GuiTextBoxId", vm.RenderGUI)
	reg("GuiDropdownBox", vm.RenderGUI)
	reg("GuiWindowBox", vm.RenderGUI)
	reg("GuiGroupBox", vm.RenderGUI)
	reg("GuiLine", vm.RenderGUI)
	reg("GuiPanel", vm.RenderGUI)
	reg("GuiWindow", vm.RenderGUI)
	reg("GuiList", vm.RenderGUI)
	reg("GuiDropdown", vm.RenderGUI)
	reg("GuiProgressBarSimple", vm.RenderGUI)
}
