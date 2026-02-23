// Package raylib: 2D shapes (rshapes).
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerShapes(v *vm.VM) {
	v.RegisterForeign("SetShapesTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("SetShapesTexture requires (textureId, srcX, srcY, srcW, srcH)")
		}
		id := toString(args[0])
		texMu.Lock()
		tex, ok := textures[id]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", id)
		}
		source := rl.Rectangle{X: toFloat32(args[1]), Y: toFloat32(args[2]), Width: toFloat32(args[3]), Height: toFloat32(args[4])}
		rl.SetShapesTexture(tex, source)
		return nil, nil
	})
	v.RegisterForeign("GetShapesTextureRectangle", func(args []interface{}) (interface{}, error) {
		rec := rl.GetShapesTextureRectangle()
		return []interface{}{float64(rec.X), float64(rec.Y), float64(rec.Width), float64(rec.Height)}, nil
	})
	v.RegisterForeign("DrawRectangle", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawRectangle requires (x, y, w, h, color or r,g,b,a)")
		}
		x, y, w, h := toInt32(args[0]), toInt32(args[1]), toInt32(args[2]), toInt32(args[3])
		var c rl.Color
		if len(args) == 5 {
			switch v := args[4].(type) {
			case int:
				c = rl.NewColor(uint8(v>>16&0xff), uint8(v>>8&0xff), uint8(v&0xff), 255)
			case float64:
				u := uint32(v)
				c = rl.NewColor(uint8(u>>16&0xff), uint8(u>>8&0xff), uint8(u&0xff), 255)
			default:
				c = rl.White
			}
		} else if len(args) >= 8 {
			c = rl.NewColor(uint8(toInt32(args[4])), uint8(toInt32(args[5])), uint8(toInt32(args[6])), uint8(toInt32(args[7])))
		} else {
			c = rl.White
		}
		rl.DrawRectangle(x, y, w, h, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCircle", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawCircle requires (centerX, centerY, radius, color)")
		}
		x, y := toInt32(args[0]), toInt32(args[1])
		radius := toFloat32(args[2])
		c := rl.White
		if len(args) >= 7 {
			c = rl.NewColor(uint8(toInt32(args[3])), uint8(toInt32(args[4])), uint8(toInt32(args[5])), uint8(toInt32(args[6])))
		}
		rl.DrawCircle(int32(x), int32(y), radius, c)
		return nil, nil
	})
	v.RegisterForeign("DrawLine", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawLine requires (x1, y1, x2, y2, color or r,g,b,a)")
		}
		x1, y1, x2, y2 := toInt32(args[0]), toInt32(args[1]), toInt32(args[2]), toInt32(args[3])
		c := rl.White
		if len(args) >= 8 {
			c = argsToColor(args, 4)
		}
		rl.DrawLine(x1, y1, x2, y2, c)
		return nil, nil
	})
	v.RegisterForeign("DrawLineV", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawLineV requires (startX, startY, endX, endY, color or r,g,b,a)")
		}
		start := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		end := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		c := rl.White
		if len(args) >= 8 {
			c = argsToColor(args, 4)
		}
		rl.DrawLineV(start, end, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCircleLines", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawCircleLines requires (centerX, centerY, radius, color)")
		}
		x, y := toInt32(args[0]), toInt32(args[1])
		radius := toFloat32(args[2])
		c := rl.White
		if len(args) >= 7 {
			c = argsToColor(args, 3)
		}
		rl.DrawCircleLines(x, y, radius, c)
		return nil, nil
	})
	v.RegisterForeign("DrawRectangleLines", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawRectangleLines requires (x, y, w, h, color)")
		}
		x, y, w, h := toInt32(args[0]), toInt32(args[1]), toInt32(args[2]), toInt32(args[3])
		c := rl.White
		if len(args) >= 8 {
			c = argsToColor(args, 4)
		}
		rl.DrawRectangleLines(x, y, w, h, c)
		return nil, nil
	})
	v.RegisterForeign("DrawTriangle", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawTriangle requires (x1,y1, x2,y2, x3,y3, color)")
		}
		v1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		v2 := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		v3 := rl.Vector2{X: toFloat32(args[4]), Y: toFloat32(args[5])}
		c := rl.White
		if len(args) >= 10 {
			c = argsToColor(args, 6)
		}
		rl.DrawTriangle(v1, v2, v3, c)
		return nil, nil
	})
	v.RegisterForeign("DrawTriangleLines", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawTriangleLines requires (x1,y1, x2,y2, x3,y3, color)")
		}
		v1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		v2 := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		v3 := rl.Vector2{X: toFloat32(args[4]), Y: toFloat32(args[5])}
		c := rl.White
		if len(args) >= 10 {
			c = argsToColor(args, 6)
		}
		rl.DrawTriangleLines(v1, v2, v3, c)
		return nil, nil
	})
	v.RegisterForeign("DrawPixel", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("DrawPixel requires (x, y, color)")
		}
		x, y := toInt32(args[0]), toInt32(args[1])
		c := rl.White
		if len(args) >= 6 {
			c = argsToColor(args, 2)
		}
		rl.DrawPixel(x, y, c)
		return nil, nil
	})
	v.RegisterForeign("DrawPoly", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("DrawPoly requires (centerX, centerY, sides, radius, rotation, color)")
		}
		cx, cy := toInt32(args[0]), toInt32(args[1])
		center := rl.Vector2{X: float32(cx), Y: float32(cy)}
		sides := toInt32(args[2])
		radius := toFloat32(args[3])
		rotation := toFloat32(args[4])
		c := rl.White
		if len(args) >= 9 {
			c = argsToColor(args, 5)
		}
		rl.DrawPoly(center, sides, radius, rotation, c)
		return nil, nil
	})
	v.RegisterForeign("DrawEllipse", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawEllipse requires (centerX, centerY, radiusH, radiusV, color)")
		}
		cx, cy := toInt32(args[0]), toInt32(args[1])
		rh, rv := toFloat32(args[2]), toFloat32(args[3])
		c := rl.White
		if len(args) >= 8 {
			c = argsToColor(args, 4)
		}
		rl.DrawEllipse(cx, cy, rh, rv, c)
		return nil, nil
	})
	v.RegisterForeign("DrawRing", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawRing requires (centerX, centerY, innerRadius, outerRadius, startAngle, endAngle, segments, color)")
		}
		center := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		innerRadius := toFloat32(args[2])
		outerRadius := toFloat32(args[3])
		startAngle := toFloat32(args[4])
		endAngle := toFloat32(args[5])
		segments := toInt32(args[6])
		if segments <= 0 {
			segments = 36
		}
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 7)
		}
		rl.DrawRing(center, innerRadius, outerRadius, startAngle, endAngle, segments, c)
		return nil, nil
	})
	v.RegisterForeign("DrawRectangleRounded", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("DrawRectangleRounded requires (x, y, w, h, roundness, segments, color)")
		}
		rec := rl.Rectangle{X: toFloat32(args[0]), Y: toFloat32(args[1]), Width: toFloat32(args[2]), Height: toFloat32(args[3])}
		roundness := toFloat32(args[4])
		segments := toInt32(args[5])
		if segments <= 0 {
			segments = 8
		}
		c := rl.White
		if len(args) >= 9 {
			c = argsToColor(args, 6)
		}
		rl.DrawRectangleRounded(rec, roundness, segments, c)
		return nil, nil
	})
	v.RegisterForeign("DrawGrid", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("DrawGrid requires (slices, spacing)")
		}
		slices := toInt32(args[0])
		spacing := toFloat32(args[1])
		rl.DrawGrid(slices, spacing)
		return nil, nil
	})
	v.RegisterForeign("DrawFPS", func(args []interface{}) (interface{}, error) {
		x, y := int32(10), int32(10)
		if len(args) >= 2 {
			x, y = toInt32(args[0]), toInt32(args[1])
		}
		rl.DrawFPS(x, y)
		return nil, nil
	})
	v.RegisterForeign("DrawLineEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("DrawLineEx requires (startX, startY, endX, endY, thick, color)")
		}
		start := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		end := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		thick := toFloat32(args[4])
		c := rl.White
		if len(args) >= 10 {
			c = argsToColor(args, 5)
		}
		rl.DrawLineEx(start, end, thick, c)
		return nil, nil
	})
	v.RegisterForeign("DrawPixelV", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("DrawPixelV requires (x, y, color)")
		}
		pos := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		c := rl.White
		if len(args) >= 7 {
			c = argsToColor(args, 2)
		}
		rl.DrawPixelV(pos, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCircleSector", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawCircleSector requires (centerX, centerY, radius, startAngle, endAngle, segments, color)")
		}
		center := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		radius := toFloat32(args[2])
		startAngle := toFloat32(args[3])
		endAngle := toFloat32(args[4])
		segments := toInt32(args[5])
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 6)
		}
		rl.DrawCircleSector(center, radius, startAngle, endAngle, segments, c)
		return nil, nil
	})
	v.RegisterForeign("DrawCircleGradient", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("DrawCircleGradient requires (centerX, centerY, radius, innerR,G,B,A, outerR,G,B,A)")
		}
		cx, cy := toInt32(args[0]), toInt32(args[1])
		radius := toFloat32(args[2])
		inner := argsToColor(args, 3)
		outer := argsToColor(args, 7)
		rl.DrawCircleGradient(cx, cy, radius, inner, outer)
		return nil, nil
	})
	v.RegisterForeign("DrawCircleV", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("DrawCircleV requires (centerX, centerY, radius, color)")
		}
		center := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		radius := toFloat32(args[2])
		c := rl.White
		if len(args) >= 7 {
			c = argsToColor(args, 3)
		}
		rl.DrawCircleV(center, radius, c)
		return nil, nil
	})
	v.RegisterForeign("DrawEllipseLines", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawEllipseLines requires (centerX, centerY, radiusH, radiusV, color)")
		}
		cx, cy := toInt32(args[0]), toInt32(args[1])
		rh, rv := toFloat32(args[2]), toFloat32(args[3])
		c := rl.White
		if len(args) >= 8 {
			c = argsToColor(args, 4)
		}
		rl.DrawEllipseLines(cx, cy, rh, rv, c)
		return nil, nil
	})
	v.RegisterForeign("DrawRingLines", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawRingLines requires (centerX, centerY, innerR, outerR, startAngle, endAngle, segments, color)")
		}
		center := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		innerRadius := toFloat32(args[2])
		outerRadius := toFloat32(args[3])
		startAngle := toFloat32(args[4])
		endAngle := toFloat32(args[5])
		segments := toInt32(args[6])
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 7)
		}
		rl.DrawRingLines(center, innerRadius, outerRadius, startAngle, endAngle, segments, c)
		return nil, nil
	})
	v.RegisterForeign("DrawRectangleV", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawRectangleV requires (posX, posY, width, height, color)")
		}
		pos := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		size := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		c := rl.White
		if len(args) >= 9 {
			c = argsToColor(args, 4)
		}
		rl.DrawRectangleV(pos, size, c)
		return nil, nil
	})
	v.RegisterForeign("DrawRectangleRec", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("DrawRectangleRec requires (x, y, width, height, color)")
		}
		rec := rl.Rectangle{X: toFloat32(args[0]), Y: toFloat32(args[1]), Width: toFloat32(args[2]), Height: toFloat32(args[3])}
		c := rl.White
		if len(args) >= 9 {
			c = argsToColor(args, 4)
		}
		rl.DrawRectangleRec(rec, c)
		return nil, nil
	})
	v.RegisterForeign("DrawRectanglePro", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("DrawRectanglePro requires (recX, recY, recW, recH, originX, originY, rotation, color)")
		}
		rec := rl.Rectangle{X: toFloat32(args[0]), Y: toFloat32(args[1]), Width: toFloat32(args[2]), Height: toFloat32(args[3])}
		origin := rl.Vector2{X: toFloat32(args[4]), Y: toFloat32(args[5])}
		rotation := toFloat32(args[6])
		c := rl.White
		if len(args) >= 11 {
			c = argsToColor(args, 7)
		}
		rl.DrawRectanglePro(rec, origin, rotation, c)
		return nil, nil
	})
	v.RegisterForeign("DrawRectangleLinesEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("DrawRectangleLinesEx requires (x, y, w, h, lineThick, color)")
		}
		rec := rl.Rectangle{X: toFloat32(args[0]), Y: toFloat32(args[1]), Width: toFloat32(args[2]), Height: toFloat32(args[3])}
		lineThick := toFloat32(args[4])
		c := rl.White
		if len(args) >= 10 {
			c = argsToColor(args, 5)
		}
		rl.DrawRectangleLinesEx(rec, lineThick, c)
		return nil, nil
	})
	v.RegisterForeign("DrawRectangleRoundedLines", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("DrawRectangleRoundedLines requires (x, y, w, h, roundness, segments, color)")
		}
		rec := rl.Rectangle{X: toFloat32(args[0]), Y: toFloat32(args[1]), Width: toFloat32(args[2]), Height: toFloat32(args[3])}
		roundness := toFloat32(args[4])
		segments := toInt32(args[5])
		c := rl.White
		if len(args) >= 9 {
			c = argsToColor(args, 6)
		}
		rl.DrawRectangleRoundedLines(rec, roundness, segments, c)
		return nil, nil
	})
	v.RegisterForeign("DrawPolyLines", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("DrawPolyLines requires (centerX, centerY, sides, radius, rotation, color)")
		}
		center := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		sides := toInt32(args[2])
		radius := toFloat32(args[3])
		rotation := toFloat32(args[4])
		c := rl.White
		if len(args) >= 9 {
			c = argsToColor(args, 5)
		}
		rl.DrawPolyLines(center, sides, radius, rotation, c)
		return nil, nil
	})
}
