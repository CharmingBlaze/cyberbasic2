// Package raylib: mouse delta, color helpers, collision, constants.
package raylib

import (
	"cyberbasic/compiler/vm"
	"fmt"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	lastRayCollision   rl.RayCollision
	lastRayCollisionMu sync.Mutex
)

func registerMisc(v *vm.VM) {
	v.RegisterForeign("GetMouseDelta", func(args []interface{}) (interface{}, error) {
		delta := rl.GetMouseDelta()
		return []interface{}{float64(delta.X), float64(delta.Y)}, nil
	})
	v.RegisterForeign("NewColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("NewColor requires (r, g, b, a)")
		}
		r, g, b, a := toUint8(args[0]), toUint8(args[1]), toUint8(args[2]), toUint8(args[3])
		return int(r)<<24 | int(g)<<16 | int(b)<<8 | int(a), nil
	})
	v.RegisterForeign("CheckCollisionRecs", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("CheckCollisionRecs requires (x1,y1,w1,h1, x2,y2,w2,h2)")
		}
		rec1 := rl.Rectangle{X: toFloat32(args[0]), Y: toFloat32(args[1]), Width: toFloat32(args[2]), Height: toFloat32(args[3])}
		rec2 := rl.Rectangle{X: toFloat32(args[4]), Y: toFloat32(args[5]), Width: toFloat32(args[6]), Height: toFloat32(args[7])}
		return rl.CheckCollisionRecs(rec1, rec2), nil
	})
	v.RegisterForeign("CheckCollisionCircles", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("CheckCollisionCircles requires (x1,y1,r1, x2,y2,r2)")
		}
		center1 := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		radius1 := toFloat32(args[2])
		center2 := rl.Vector2{X: toFloat32(args[3]), Y: toFloat32(args[4])}
		radius2 := toFloat32(args[5])
		return rl.CheckCollisionCircles(center1, radius1, center2, radius2), nil
	})
	v.RegisterForeign("CheckCollisionCircleRec", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("CheckCollisionCircleRec requires (centerX, centerY, radius, recX, recY, recW, recH)")
		}
		center := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		radius := toFloat32(args[2])
		rec := rl.Rectangle{X: toFloat32(args[3]), Y: toFloat32(args[4]), Width: toFloat32(args[5]), Height: toFloat32(args[6])}
		return rl.CheckCollisionCircleRec(center, radius, rec), nil
	})
	v.RegisterForeign("CheckCollisionPointRec", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("CheckCollisionPointRec requires (pointX, pointY, recX, recY, recW, recH)")
		}
		point := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		rec := rl.Rectangle{X: toFloat32(args[2]), Y: toFloat32(args[3]), Width: toFloat32(args[4]), Height: toFloat32(args[5])}
		return rl.CheckCollisionPointRec(point, rec), nil
	})
	v.RegisterForeign("CheckCollisionPointCircle", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("CheckCollisionPointCircle requires (pointX, pointY, centerX, centerY, radius)")
		}
		point := rl.Vector2{X: toFloat32(args[0]), Y: toFloat32(args[1])}
		center := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		radius := toFloat32(args[4])
		return rl.CheckCollisionPointCircle(point, center, radius), nil
	})
	v.RegisterForeign("GetCollisionRec", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("GetCollisionRec requires (x1,y1,w1,h1, x2,y2,w2,h2)")
		}
		rec1 := rl.Rectangle{X: toFloat32(args[0]), Y: toFloat32(args[1]), Width: toFloat32(args[2]), Height: toFloat32(args[3])}
		rec2 := rl.Rectangle{X: toFloat32(args[4]), Y: toFloat32(args[5]), Width: toFloat32(args[6]), Height: toFloat32(args[7])}
		coll := rl.GetCollisionRec(rec1, rec2)
		return []interface{}{float64(coll.X), float64(coll.Y), float64(coll.Width), float64(coll.Height)}, nil
	})
	// 3D collision
	v.RegisterForeign("CheckCollisionSpheres", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("CheckCollisionSpheres requires (center1X,1Y,1Z, radius1, center2X,2Y,2Z, radius2)")
		}
		center1 := rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])}
		radius1 := toFloat32(args[3])
		center2 := rl.Vector3{X: toFloat32(args[4]), Y: toFloat32(args[5]), Z: toFloat32(args[6])}
		radius2 := toFloat32(args[7])
		return rl.CheckCollisionSpheres(center1, radius1, center2, radius2), nil
	})
	v.RegisterForeign("CheckCollisionBoxes", func(args []interface{}) (interface{}, error) {
		if len(args) < 12 {
			return nil, fmt.Errorf("CheckCollisionBoxes requires (box1MinX,Y,Z, box1MaxX,Y,Z, box2MinX,Y,Z, box2MaxX,Y,Z)")
		}
		box1 := rl.BoundingBox{
			Min: rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])},
			Max: rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])},
		}
		box2 := rl.BoundingBox{
			Min: rl.Vector3{X: toFloat32(args[6]), Y: toFloat32(args[7]), Z: toFloat32(args[8])},
			Max: rl.Vector3{X: toFloat32(args[9]), Y: toFloat32(args[10]), Z: toFloat32(args[11])},
		}
		return rl.CheckCollisionBoxes(box1, box2), nil
	})
	v.RegisterForeign("CheckCollisionBoxSphere", func(args []interface{}) (interface{}, error) {
		if len(args) < 10 {
			return nil, fmt.Errorf("CheckCollisionBoxSphere requires (boxMinX,Y,Z, boxMaxX,Y,Z, centerX,Y,Z, radius)")
		}
		box := rl.BoundingBox{
			Min: rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])},
			Max: rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])},
		}
		center := rl.Vector3{X: toFloat32(args[6]), Y: toFloat32(args[7]), Z: toFloat32(args[8])}
		radius := toFloat32(args[9])
		return rl.CheckCollisionBoxSphere(box, center, radius), nil
	})
	// Ray collision: store result in lastRayCollision, return hit (1/0)
	v.RegisterForeign("GetRayCollisionSphere", func(args []interface{}) (interface{}, error) {
		if len(args) < 10 {
			return nil, fmt.Errorf("GetRayCollisionSphere requires (rayPosX,Y,Z, rayDirX,Y,Z, centerX,Y,Z, radius)")
		}
		ray := rl.Ray{
			Position:  rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])},
			Direction: rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])},
		}
		center := rl.Vector3{X: toFloat32(args[6]), Y: toFloat32(args[7]), Z: toFloat32(args[8])}
		radius := toFloat32(args[9])
		coll := rl.GetRayCollisionSphere(ray, center, radius)
		lastRayCollisionMu.Lock()
		lastRayCollision = coll
		lastRayCollisionMu.Unlock()
		if coll.Hit {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("GetRayCollisionBox", func(args []interface{}) (interface{}, error) {
		if len(args) < 13 {
			return nil, fmt.Errorf("GetRayCollisionBox requires (rayPosX,Y,Z, rayDirX,Y,Z, boxMinX,Y,Z, boxMaxX,Y,Z)")
		}
		ray := rl.Ray{
			Position:  rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])},
			Direction: rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])},
		}
		box := rl.BoundingBox{
			Min: rl.Vector3{X: toFloat32(args[6]), Y: toFloat32(args[7]), Z: toFloat32(args[8])},
			Max: rl.Vector3{X: toFloat32(args[9]), Y: toFloat32(args[10]), Z: toFloat32(args[11])},
		}
		coll := rl.GetRayCollisionBox(ray, box)
		lastRayCollisionMu.Lock()
		lastRayCollision = coll
		lastRayCollisionMu.Unlock()
		if coll.Hit {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("GetRayCollisionTriangle", func(args []interface{}) (interface{}, error) {
		if len(args) < 12 {
			return nil, fmt.Errorf("GetRayCollisionTriangle requires (rayPosX,Y,Z, rayDirX,Y,Z, p1x,p1y,p1z, p2x,p2y,p2z, p3x,p3y,p3z)")
		}
		ray := rl.Ray{
			Position:  rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])},
			Direction: rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])},
		}
		p1 := rl.Vector3{X: toFloat32(args[6]), Y: toFloat32(args[7]), Z: toFloat32(args[8])}
		p2 := rl.Vector3{X: toFloat32(args[9]), Y: toFloat32(args[10]), Z: toFloat32(args[11])}
		p3 := rl.Vector3{X: toFloat32(args[12]), Y: toFloat32(args[13]), Z: toFloat32(args[14])}
		coll := rl.GetRayCollisionTriangle(ray, p1, p2, p3)
		lastRayCollisionMu.Lock()
		lastRayCollision = coll
		lastRayCollisionMu.Unlock()
		if coll.Hit {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("GetRayCollisionQuad", func(args []interface{}) (interface{}, error) {
		if len(args) < 18 {
			return nil, fmt.Errorf("GetRayCollisionQuad requires (rayPosX,Y,Z, rayDirX,Y,Z, p1x,p1y,p1z, p2x,p2y,p2z, p3x,p3y,p3z, p4x,p4y,p4z)")
		}
		ray := rl.Ray{
			Position:  rl.Vector3{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2])},
			Direction: rl.Vector3{X: toFloat32(args[3]), Y: toFloat32(args[4]), Z: toFloat32(args[5])},
		}
		p1 := rl.Vector3{X: toFloat32(args[6]), Y: toFloat32(args[7]), Z: toFloat32(args[8])}
		p2 := rl.Vector3{X: toFloat32(args[9]), Y: toFloat32(args[10]), Z: toFloat32(args[11])}
		p3 := rl.Vector3{X: toFloat32(args[12]), Y: toFloat32(args[13]), Z: toFloat32(args[14])}
		p4 := rl.Vector3{X: toFloat32(args[15]), Y: toFloat32(args[16]), Z: toFloat32(args[17])}
		coll := rl.GetRayCollisionQuad(ray, p1, p2, p3, p4)
		lastRayCollisionMu.Lock()
		lastRayCollision = coll
		lastRayCollisionMu.Unlock()
		if coll.Hit {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("GetRayCollisionPointX", func(args []interface{}) (interface{}, error) {
		lastRayCollisionMu.Lock()
		defer lastRayCollisionMu.Unlock()
		return float64(lastRayCollision.Point.X), nil
	})
	v.RegisterForeign("GetRayCollisionPointY", func(args []interface{}) (interface{}, error) {
		lastRayCollisionMu.Lock()
		defer lastRayCollisionMu.Unlock()
		return float64(lastRayCollision.Point.Y), nil
	})
	v.RegisterForeign("GetRayCollisionPointZ", func(args []interface{}) (interface{}, error) {
		lastRayCollisionMu.Lock()
		defer lastRayCollisionMu.Unlock()
		return float64(lastRayCollision.Point.Z), nil
	})
	v.RegisterForeign("GetRayCollisionNormalX", func(args []interface{}) (interface{}, error) {
		lastRayCollisionMu.Lock()
		defer lastRayCollisionMu.Unlock()
		return float64(lastRayCollision.Normal.X), nil
	})
	v.RegisterForeign("GetRayCollisionNormalY", func(args []interface{}) (interface{}, error) {
		lastRayCollisionMu.Lock()
		defer lastRayCollisionMu.Unlock()
		return float64(lastRayCollision.Normal.Y), nil
	})
	v.RegisterForeign("GetRayCollisionNormalZ", func(args []interface{}) (interface{}, error) {
		lastRayCollisionMu.Lock()
		defer lastRayCollisionMu.Unlock()
		return float64(lastRayCollision.Normal.Z), nil
	})
	v.RegisterForeign("GetRayCollisionDistance", func(args []interface{}) (interface{}, error) {
		lastRayCollisionMu.Lock()
		defer lastRayCollisionMu.Unlock()
		return float64(lastRayCollision.Distance), nil
	})
	v.RegisterForeign("Fade", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("Fade requires (r, g, b, a, alpha)")
		}
		c := rl.NewColor(toUint8(args[0]), toUint8(args[1]), toUint8(args[2]), toUint8(args[3]))
		out := rl.Fade(c, toFloat32(args[4]))
		return int(out.R)<<24|int(out.G)<<16|int(out.B)<<8|int(out.A), nil
	})
	v.RegisterForeign("ColorAlpha", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ColorAlpha requires (r, g, b, a, alpha)")
		}
		c := rl.NewColor(toUint8(args[0]), toUint8(args[1]), toUint8(args[2]), toUint8(args[3]))
		out := rl.ColorAlpha(c, toFloat32(args[4]))
		return int(out.R)<<24|int(out.G)<<16|int(out.B)<<8|int(out.A), nil
	})
	v.RegisterForeign("ColorToInt", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ColorToInt requires (r, g, b, a)")
		}
		c := rl.NewColor(toUint8(args[0]), toUint8(args[1]), toUint8(args[2]), toUint8(args[3]))
		return int(rl.ColorToInt(c)), nil
	})
	v.RegisterForeign("GetColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetColor requires (hexValue)")
		}
		hexVal := toInt32(args[0])
		if h, ok := args[0].(int); ok {
			c := rl.GetColor(uint(h))
			return []interface{}{int(c.R), int(c.G), int(c.B), int(c.A)}, nil
		}
		c := rl.GetColor(uint(hexVal))
		return []interface{}{int(c.R), int(c.G), int(c.B), int(c.A)}, nil
	})
	// Color constants return packed int (R<<16|G<<8|B) for use with Draw* and text functions.
	v.RegisterForeign("White", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.White), nil })
	v.RegisterForeign("Black", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Black), nil })
	v.RegisterForeign("LightGray", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.LightGray), nil })
	v.RegisterForeign("Gray", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Gray), nil })
	v.RegisterForeign("DarkGray", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.DarkGray), nil })
	v.RegisterForeign("Yellow", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Yellow), nil })
	v.RegisterForeign("Gold", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Gold), nil })
	v.RegisterForeign("Orange", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Orange), nil })
	v.RegisterForeign("Pink", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Pink), nil })
	v.RegisterForeign("Red", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Red), nil })
	v.RegisterForeign("Maroon", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Maroon), nil })
	v.RegisterForeign("Green", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Green), nil })
	v.RegisterForeign("Lime", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Lime), nil })
	v.RegisterForeign("DarkGreen", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.DarkGreen), nil })
	v.RegisterForeign("SkyBlue", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.SkyBlue), nil })
	v.RegisterForeign("Blue", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Blue), nil })
	v.RegisterForeign("DarkBlue", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.DarkBlue), nil })
	v.RegisterForeign("Purple", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Purple), nil })
	v.RegisterForeign("Violet", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Violet), nil })
	v.RegisterForeign("DarkPurple", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.DarkPurple), nil })
	v.RegisterForeign("Beige", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Beige), nil })
	v.RegisterForeign("Brown", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Brown), nil })
	v.RegisterForeign("DarkBrown", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.DarkBrown), nil })
	v.RegisterForeign("Magenta", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Magenta), nil })
	v.RegisterForeign("RayWhite", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.RayWhite), nil })
	v.RegisterForeign("Blank", func(args []interface{}) (interface{}, error) { return colorToPacked(rl.Blank), nil })

	// Color utilities (match C API)
	v.RegisterForeign("ColorIsEqual", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return false, nil
		}
		c1 := argsToColor(args, 0)
		c2 := argsToColor(args, 4)
		return c1.R == c2.R && c1.G == c2.G && c1.B == c2.B && c1.A == c2.A, nil
	})
	v.RegisterForeign("ColorNormalize", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ColorNormalize requires (r, g, b, a)")
		}
		c := argsToColor(args, 0)
		v4 := rl.ColorNormalize(c)
		return []interface{}{float64(v4.X), float64(v4.Y), float64(v4.Z), float64(v4.W)}, nil
	})
	v.RegisterForeign("ColorFromNormalized", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ColorFromNormalized requires (x, y, z, w)")
		}
		v4 := rl.Vector4{X: toFloat32(args[0]), Y: toFloat32(args[1]), Z: toFloat32(args[2]), W: toFloat32(args[3])}
		c := rl.ColorFromNormalized(v4)
		return int(c.R)<<24 | int(c.G)<<16 | int(c.B)<<8 | int(c.A), nil
	})
	v.RegisterForeign("ColorToHSV", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ColorToHSV requires (r, g, b, a)")
		}
		c := argsToColor(args, 0)
		v3 := rl.ColorToHSV(c)
		return []interface{}{float64(v3.X), float64(v3.Y), float64(v3.Z)}, nil
	})
	v.RegisterForeign("ColorFromHSV", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ColorFromHSV requires (hue, saturation, value)")
		}
		c := rl.ColorFromHSV(toFloat32(args[0]), toFloat32(args[1]), toFloat32(args[2]))
		return int(c.R)<<24 | int(c.G)<<16 | int(c.B)<<8 | int(c.A), nil
	})
	v.RegisterForeign("ColorTint", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("ColorTint requires (r, g, b, a, tintR, tintG, tintB, tintA)")
		}
		c := argsToColor(args, 0)
		tint := argsToColor(args, 4)
		out := rl.ColorTint(c, tint)
		return int(out.R)<<24 | int(out.G)<<16 | int(out.B)<<8 | int(out.A), nil
	})
	v.RegisterForeign("ColorBrightness", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ColorBrightness requires (r, g, b, a, factor)")
		}
		c := argsToColor(args, 0)
		out := rl.ColorBrightness(c, toFloat32(args[4]))
		return int(out.R)<<24 | int(out.G)<<16 | int(out.B)<<8 | int(out.A), nil
	})
	v.RegisterForeign("ColorContrast", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ColorContrast requires (r, g, b, a, contrast)")
		}
		c := argsToColor(args, 0)
		out := rl.ColorContrast(c, toFloat32(args[4]))
		return int(out.R)<<24 | int(out.G)<<16 | int(out.B)<<8 | int(out.A), nil
	})
	v.RegisterForeign("ColorAlphaBlend", func(args []interface{}) (interface{}, error) {
		if len(args) < 12 {
			return nil, fmt.Errorf("ColorAlphaBlend requires (dstR,g,b,a, srcR,g,b,a, tintR,g,b,a)")
		}
		dst := argsToColor(args, 0)
		src := argsToColor(args, 4)
		tint := argsToColor(args, 8)
		out := rl.ColorAlphaBlend(src, dst, tint)
		return int(out.R)<<24 | int(out.G)<<16 | int(out.B)<<8 | int(out.A), nil
	})
	v.RegisterForeign("ColorLerp", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("ColorLerp requires (r1,g1,b1,a1, r2,g2,b2,a2, factor)")
		}
		c1 := argsToColor(args, 0)
		c2 := argsToColor(args, 4)
		factor := toFloat32(args[8])
		out := rl.ColorLerp(c1, c2, factor)
		return int(out.R)<<24 | int(out.G)<<16 | int(out.B)<<8 | int(out.A), nil
	})
	v.RegisterForeign("GetPixelDataSize", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GetPixelDataSize requires (width, height, format)")
		}
		return int(rl.GetPixelDataSize(toInt32(args[0]), toInt32(args[1]), toInt32(args[2]))), nil
	})
}
