// Package raylib: image load, gen, manipulation, drawing (CPU).
package raylib

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"cyberbasic/compiler/vm"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerImages(v *vm.VM) {
	// --- Image loading ---
	v.RegisterForeign("LoadImage", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadImage requires (fileName)")
		}
		img := rl.LoadImage(toString(args[0]))
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadImageRaw", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("LoadImageRaw requires (fileName, width, height, format, headerSize)")
		}
		img := rl.LoadImageRaw(toString(args[0]), toInt32(args[1]), toInt32(args[2]), rl.PixelFormat(toInt32(args[3])), toInt32(args[4]))
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadImageAnim", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadImageAnim requires (fileName)")
		}
		var frames int32
		img := rl.LoadImageAnim(toString(args[0]), &frames)
		if img == nil {
			return "", nil
		}
		lastLoadImageAnimFrames = frames
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GetLoadImageAnimFrames", func(args []interface{}) (interface{}, error) {
		return int(lastLoadImageAnimFrames), nil
	})
	v.RegisterForeign("LoadImageAnimFromMemory", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("LoadImageAnimFromMemory requires (fileType, data, dataSize)")
		}
		var data []byte
		switch d := args[1].(type) {
		case string:
			data = []byte(d)
		case []byte:
			data = d
		default:
			return nil, fmt.Errorf("data must be string or []byte")
		}
		dataSize := toInt32(args[2])
		if int(dataSize) < len(data) {
			data = data[:dataSize]
		}
		var frames int32
		img := rl.LoadImageAnimFromMemory(toString(args[0]), data, dataSize, &frames)
		if img == nil {
			return "", nil
		}
		lastLoadImageAnimFrames = frames
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadImageFromMemory", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("LoadImageFromMemory requires (fileType, data, dataSize)")
		}
		var data []byte
		switch d := args[1].(type) {
		case string:
			data = []byte(d)
		case []byte:
			data = d
		default:
			return nil, fmt.Errorf("data must be string or []byte")
		}
		dataSize := toInt32(args[2])
		if int(dataSize) < len(data) {
			data = data[:dataSize]
		}
		img := rl.LoadImageFromMemory(toString(args[0]), data, dataSize)
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadImageFromTexture", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadImageFromTexture requires (textureId)")
		}
		texMu.Lock()
		tex, ok := textures[toString(args[0])]
		texMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown texture id: %s", toString(args[0]))
		}
		img := rl.LoadImageFromTexture(tex)
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("LoadImageFromScreen", func(args []interface{}) (interface{}, error) {
		img := rl.LoadImageFromScreen()
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("IsImageValid", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return false, nil
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		return ok && img != nil && rl.IsImageValid(img), nil
	})
	v.RegisterForeign("UnloadImage", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("UnloadImage requires (imageId)")
		}
		id := toString(args[0])
		imageMu.Lock()
		img, ok := images[id]
		delete(images, id)
		imageMu.Unlock()
		if ok && img != nil {
			rl.UnloadImage(img)
		}
		return nil, nil
	})
	v.RegisterForeign("ExportImage", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ExportImage requires (imageId, fileName)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return false, nil
		}
		return rl.ExportImage(*img, toString(args[1])), nil
	})
	v.RegisterForeign("ExportImageToMemory", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ExportImageToMemory requires (imageId, fileType)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return "", nil
		}
		data := rl.ExportImageToMemory(*img, toString(args[1]))
		if data == nil {
			return "", nil
		}
		return base64.StdEncoding.EncodeToString(data), nil
	})
	// ExportImageAsCode: export image pixel data as C header (raylib-go has no native API; we generate .h from LoadImageColors).
	v.RegisterForeign("ExportImageAsCode", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ExportImageAsCode requires (imageId, fileName)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return false, nil
		}
		cols := rl.LoadImageColors(img)
		if cols == nil {
			return false, nil
		}
		defer rl.UnloadImageColors(cols)
		w, h := img.Width, img.Height
		var b strings.Builder
		b.WriteString("// Exported by CyberBasic ExportImageAsCode (RGBA)\n")
		b.WriteString("#ifndef IMAGE_EXPORT_H\n#define IMAGE_EXPORT_H\n\n")
		fmt.Fprintf(&b, "static const int IMAGE_WIDTH = %d;\n", w)
		fmt.Fprintf(&b, "static const int IMAGE_HEIGHT = %d;\n", h)
		fmt.Fprintf(&b, "static const unsigned char IMAGE_DATA[] = {\n")
		for i, c := range cols {
			if i > 0 {
				b.WriteByte(',')
			}
			if i%16 == 0 {
				b.WriteString("\n    ")
			}
			fmt.Fprintf(&b, "%d,%d,%d,%d", c.R, c.G, c.B, c.A)
		}
		b.WriteString("\n};\n\n#endif\n")
		if err := os.WriteFile(toString(args[1]), []byte(b.String()), 0644); err != nil {
			return false, err
		}
		return true, nil
	})

	// --- Image generation ---
	v.RegisterForeign("GenImageColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("GenImageColor requires (width, height, r, g, b, a)")
		}
		c := argsToColor(args, 2)
		img := rl.GenImageColor(int(toInt32(args[0])), int(toInt32(args[1])), c)
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenImageGradientLinear", func(args []interface{}) (interface{}, error) {
		if len(args) < 10 {
			return nil, fmt.Errorf("GenImageGradientLinear requires (width, height, direction, startR,g,b,a, endR,g,b,a)")
		}
		start := argsToColor(args, 3)
		end := argsToColor(args, 7)
		img := rl.GenImageGradientLinear(int(toInt32(args[0])), int(toInt32(args[1])), int(toInt32(args[2])), start, end)
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenImageGradientRadial", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("GenImageGradientRadial requires (width, height, density, innerR,g,b,a, outerR,g,b,a)")
		}
		inner := argsToColor(args, 3)
		outer := argsToColor(args, 7)
		img := rl.GenImageGradientRadial(int(toInt32(args[0])), int(toInt32(args[1])), toFloat32(args[2]), inner, outer)
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenImageGradientSquare", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("GenImageGradientSquare requires (width, height, density, innerR,g,b,a, outerR,g,b,a)")
		}
		inner := argsToColor(args, 3)
		outer := argsToColor(args, 7)
		img := rl.GenImageGradientSquare(int(toInt32(args[0])), int(toInt32(args[1])), toFloat32(args[2]), inner, outer)
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenImageChecked", func(args []interface{}) (interface{}, error) {
		if len(args) < 10 {
			return nil, fmt.Errorf("GenImageChecked requires (width, height, checksX, checksY, col1r,g,b,a, col2r,g,b,a)")
		}
		col1 := argsToColor(args, 4)
		col2 := argsToColor(args, 8)
		img := rl.GenImageChecked(int(toInt32(args[0])), int(toInt32(args[1])), int(toInt32(args[2])), int(toInt32(args[3])), col1, col2)
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenImageWhiteNoise", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GenImageWhiteNoise requires (width, height, factor)")
		}
		img := rl.GenImageWhiteNoise(int(toInt32(args[0])), int(toInt32(args[1])), toFloat32(args[2]))
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenImagePerlinNoise", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("GenImagePerlinNoise requires (width, height, offsetX, offsetY, scale)")
		}
		img := rl.GenImagePerlinNoise(int(toInt32(args[0])), int(toInt32(args[1])), int(toInt32(args[2])), int(toInt32(args[3])), toFloat32(args[4]))
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenImageCellular", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GenImageCellular requires (width, height, tileSize)")
		}
		img := rl.GenImageCellular(int(toInt32(args[0])), int(toInt32(args[1])), int(toInt32(args[2])))
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("GenImageText", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GenImageText requires (width, height, text)")
		}
		img := rl.GenImageText(int(toInt32(args[0])), int(toInt32(args[1])), toString(args[2]))
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})

	// --- Image manipulation (return new image) ---
	v.RegisterForeign("ImageCopy", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ImageCopy requires (imageId)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		copyImg := rl.ImageCopy(img)
		if copyImg == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = copyImg
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("ImageFromImage", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ImageFromImage requires (imageId, srcX, srcY, srcW, srcH)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rec := rl.Rectangle{X: toFloat32(args[1]), Y: toFloat32(args[2]), Width: toFloat32(args[3]), Height: toFloat32(args[4])}
		result := rl.ImageFromImage(*img, rec)
		dup := &result
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = dup
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("ImageFromChannel", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ImageFromChannel requires (imageId, channel)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		result := rl.ImageFromChannel(*img, toInt32(args[1]))
		dup := &result
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = dup
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("ImageText", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("ImageText requires (text, fontSize, r, g, b, a)")
		}
		c := argsToColor(args, 2)
		img := rl.ImageText(toString(args[0]), toInt32(args[1]), c)
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("ImageTextEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("ImageTextEx requires (fontId, text, fontSize, spacing, r, g, b, a)")
		}
		fontMu.Lock()
		font, ok := fonts[toString(args[0])]
		fontMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown font id: %s", toString(args[0]))
		}
		c := argsToColor(args, 4)
		img := rl.ImageTextEx(font, toString(args[1]), toFloat32(args[2]), toFloat32(args[3]), c)
		if img == nil {
			return "", nil
		}
		imageMu.Lock()
		imageCounter++
		id := fmt.Sprintf("img_%d", imageCounter)
		images[id] = img
		imageMu.Unlock()
		return id, nil
	})

	// --- Image manipulation (in-place) ---
	v.RegisterForeign("ImageFormat", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ImageFormat requires (imageId, newFormat)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageFormat(img, rl.PixelFormat(toInt32(args[1])))
		return nil, nil
	})
	v.RegisterForeign("ImageToPOT", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ImageToPOT requires (imageId, fillR, fillG, fillB, fillA)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageToPOT(img, argsToColor(args, 1))
		return nil, nil
	})
	v.RegisterForeign("ImageCrop", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ImageCrop requires (imageId, x, y, w, h)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rec := rl.Rectangle{X: toFloat32(args[1]), Y: toFloat32(args[2]), Width: toFloat32(args[3]), Height: toFloat32(args[4])}
		rl.ImageCrop(img, rec)
		return nil, nil
	})
	v.RegisterForeign("ImageAlphaCrop", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ImageAlphaCrop requires (imageId, threshold)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageAlphaCrop(img, toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("ImageAlphaClear", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("ImageAlphaClear requires (imageId, r, g, b, a, threshold)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageAlphaClear(img, argsToColor(args, 1), toFloat32(args[5]))
		return nil, nil
	})
	v.RegisterForeign("ImageAlphaMask", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ImageAlphaMask requires (imageId, alphaMaskImageId)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		mask, ok2 := images[toString(args[1])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		if !ok2 || mask == nil {
			return nil, fmt.Errorf("unknown alpha mask image id: %s", toString(args[1]))
		}
		rl.ImageAlphaMask(img, mask)
		return nil, nil
	})
	v.RegisterForeign("ImageAlphaPremultiply", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ImageAlphaPremultiply requires (imageId)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageAlphaPremultiply(img)
		return nil, nil
	})
	v.RegisterForeign("ImageBlurGaussian", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ImageBlurGaussian requires (imageId, blurSize)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageBlurGaussian(img, toInt32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("ImageKernelConvolution", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ImageKernelConvolution requires (imageId, kernelSize, ...floats)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		kernelSize := int(toInt32(args[1]))
		need := kernelSize * kernelSize
		if len(args) < 2+need {
			return nil, fmt.Errorf("ImageKernelConvolution needs %d kernel floats", need)
		}
		kernel := make([]float32, need)
		for i := 0; i < need; i++ {
			kernel[i] = toFloat32(args[2+i])
		}
		rl.ImageKernelConvolution(img, kernel)
		return nil, nil
	})
	v.RegisterForeign("ImageResize", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ImageResize requires (imageId, newWidth, newHeight)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageResize(img, toInt32(args[1]), toInt32(args[2]))
		return nil, nil
	})
	v.RegisterForeign("ImageResizeNN", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("ImageResizeNN requires (imageId, newWidth, newHeight)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageResizeNN(img, toInt32(args[1]), toInt32(args[2]))
		return nil, nil
	})
	v.RegisterForeign("ImageResizeCanvas", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("ImageResizeCanvas requires (imageId, newWidth, newHeight, offsetX, offsetY, fillR, fillG, fillB, fillA)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageResizeCanvas(img, toInt32(args[1]), toInt32(args[2]), toInt32(args[3]), toInt32(args[4]), argsToColor(args, 5))
		return nil, nil
	})
	v.RegisterForeign("ImageMipmaps", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ImageMipmaps requires (imageId)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageMipmaps(img)
		return nil, nil
	})
	v.RegisterForeign("ImageDither", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ImageDither requires (imageId, rBpp, gBpp, bBpp, aBpp)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageDither(img, toInt32(args[1]), toInt32(args[2]), toInt32(args[3]), toInt32(args[4]))
		return nil, nil
	})
	v.RegisterForeign("ImageFlipVertical", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ImageFlipVertical requires (imageId)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageFlipVertical(img)
		return nil, nil
	})
	v.RegisterForeign("ImageFlipHorizontal", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ImageFlipHorizontal requires (imageId)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageFlipHorizontal(img)
		return nil, nil
	})
	v.RegisterForeign("ImageRotate", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ImageRotate requires (imageId, degrees)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageRotate(img, toInt32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("ImageRotateCW", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ImageRotateCW requires (imageId)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageRotateCW(img)
		return nil, nil
	})
	v.RegisterForeign("ImageRotateCCW", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ImageRotateCCW requires (imageId)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageRotateCCW(img)
		return nil, nil
	})
	v.RegisterForeign("ImageColorTint", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ImageColorTint requires (imageId, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageColorTint(img, argsToColor(args, 1))
		return nil, nil
	})
	v.RegisterForeign("ImageColorInvert", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ImageColorInvert requires (imageId)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageColorInvert(img)
		return nil, nil
	})
	v.RegisterForeign("ImageColorGrayscale", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ImageColorGrayscale requires (imageId)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageColorGrayscale(img)
		return nil, nil
	})
	v.RegisterForeign("ImageColorContrast", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ImageColorContrast requires (imageId, contrast)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageColorContrast(img, toFloat32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("ImageColorBrightness", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ImageColorBrightness requires (imageId, brightness)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageColorBrightness(img, toInt32(args[1]))
		return nil, nil
	})
	v.RegisterForeign("ImageColorReplace", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("ImageColorReplace requires (imageId, colorR,g,b,a, replaceR,g,b,a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageColorReplace(img, argsToColor(args, 1), argsToColor(args, 5))
		return nil, nil
	})

	// LoadImageColors / UnloadImageColors / GetLoadedImageColor
	v.RegisterForeign("LoadImageColors", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LoadImageColors requires (imageId)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		cols := rl.LoadImageColors(img)
		lastImageColorsMu.Lock()
		lastImageColors = cols
		lastImageColorsMu.Unlock()
		return len(cols), nil
	})
	v.RegisterForeign("UnloadImageColors", func(args []interface{}) (interface{}, error) {
		lastImageColorsMu.Lock()
		if len(lastImageColors) > 0 {
			rl.UnloadImageColors(lastImageColors)
			lastImageColors = nil
		}
		lastImageColorsMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("GetLoadedImageColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, nil
		}
		idx := int(toInt32(args[0]))
		lastImageColorsMu.Lock()
		defer lastImageColorsMu.Unlock()
		if idx < 0 || idx >= len(lastImageColors) {
			return nil, nil
		}
		c := lastImageColors[idx]
		return []interface{}{int(c.R), int(c.G), int(c.B), int(c.A)}, nil
	})

	// GetImageColor (single pixel)
	v.RegisterForeign("GetImageColor", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("GetImageColor requires (imageId, x, y)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		c := rl.GetImageColor(*img, toInt32(args[1]), toInt32(args[2]))
		return []interface{}{int(c.R), int(c.G), int(c.B), int(c.A)}, nil
	})

	// Image drawing
	v.RegisterForeign("ImageClearBackground", func(args []interface{}) (interface{}, error) {
		if len(args) < 5 {
			return nil, fmt.Errorf("ImageClearBackground requires (imageId, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageClearBackground(img, argsToColor(args, 1))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawPixel", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("ImageDrawPixel requires (imageId, posX, posY, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageDrawPixel(img, toInt32(args[1]), toInt32(args[2]), argsToColor(args, 3))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawPixelV", func(args []interface{}) (interface{}, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("ImageDrawPixelV requires (imageId, x, y, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageDrawPixelV(img, rl.Vector2{X: toFloat32(args[1]), Y: toFloat32(args[2])}, argsToColor(args, 3))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawLine", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("ImageDrawLine requires (imageId, startX, startY, endX, endY, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageDrawLine(img, toInt32(args[1]), toInt32(args[2]), toInt32(args[3]), toInt32(args[4]), argsToColor(args, 5))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawLineV", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("ImageDrawLineV requires (imageId, startX, startY, endX, endY, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageDrawLineV(img, rl.Vector2{X: toFloat32(args[1]), Y: toFloat32(args[2])}, rl.Vector2{X: toFloat32(args[3]), Y: toFloat32(args[4])}, argsToColor(args, 5))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawLineEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("ImageDrawLineEx requires (imageId, x1, y1, x2, y2, thick, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageDrawLineEx(img, rl.Vector2{X: toFloat32(args[1]), Y: toFloat32(args[2])}, rl.Vector2{X: toFloat32(args[3]), Y: toFloat32(args[4])}, toInt32(args[5]), argsToColor(args, 6))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawCircle", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("ImageDrawCircle requires (imageId, centerX, centerY, radius, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageDrawCircle(img, toInt32(args[1]), toInt32(args[2]), toInt32(args[3]), argsToColor(args, 4))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawCircleV", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("ImageDrawCircleV requires (imageId, centerX, centerY, radius, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageDrawCircleV(img, rl.Vector2{X: toFloat32(args[1]), Y: toFloat32(args[2])}, toInt32(args[3]), argsToColor(args, 4))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawCircleLines", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("ImageDrawCircleLines requires (imageId, centerX, centerY, radius, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageDrawCircleLines(img, toInt32(args[1]), toInt32(args[2]), toInt32(args[3]), argsToColor(args, 4))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawCircleLinesV", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("ImageDrawCircleLinesV requires (imageId, centerX, centerY, radius, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageDrawCircleLinesV(img, rl.Vector2{X: toFloat32(args[1]), Y: toFloat32(args[2])}, toInt32(args[3]), argsToColor(args, 4))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawRectangle", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("ImageDrawRectangle requires (imageId, x, y, width, height, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageDrawRectangle(img, toInt32(args[1]), toInt32(args[2]), toInt32(args[3]), toInt32(args[4]), argsToColor(args, 5))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawRectangleV", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("ImageDrawRectangleV requires (imageId, posX, posY, sizeX, sizeY, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rl.ImageDrawRectangleV(img, rl.Vector2{X: toFloat32(args[1]), Y: toFloat32(args[2])}, rl.Vector2{X: toFloat32(args[3]), Y: toFloat32(args[4])}, argsToColor(args, 5))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawRectangleRec", func(args []interface{}) (interface{}, error) {
		if len(args) < 8 {
			return nil, fmt.Errorf("ImageDrawRectangleRec requires (imageId, x, y, w, h, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rec := rl.Rectangle{X: toFloat32(args[1]), Y: toFloat32(args[2]), Width: toFloat32(args[3]), Height: toFloat32(args[4])}
		rl.ImageDrawRectangleRec(img, rec, argsToColor(args, 5))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawRectangleLines", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("ImageDrawRectangleLines requires (imageId, x, y, w, h, thick, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		rec := rl.Rectangle{X: toFloat32(args[1]), Y: toFloat32(args[2]), Width: toFloat32(args[3]), Height: toFloat32(args[4])}
		rl.ImageDrawRectangleLines(img, rec, int(toInt32(args[5])), argsToColor(args, 6))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawTriangle", func(args []interface{}) (interface{}, error) {
		if len(args) < 10 {
			return nil, fmt.Errorf("ImageDrawTriangle requires (imageId, x1,y1, x2,y2, x3,y3, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		v1 := rl.Vector2{X: toFloat32(args[1]), Y: toFloat32(args[2])}
		v2 := rl.Vector2{X: toFloat32(args[3]), Y: toFloat32(args[4])}
		v3 := rl.Vector2{X: toFloat32(args[5]), Y: toFloat32(args[6])}
		rl.ImageDrawTriangle(img, v1, v2, v3, argsToColor(args, 7))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawTriangleEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 16 {
			return nil, fmt.Errorf("ImageDrawTriangleEx requires (imageId, x1,y1, x2,y2, x3,y3, c1r,g,b,a, c2r,g,b,a, c3r,g,b,a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		v1 := rl.Vector2{X: toFloat32(args[1]), Y: toFloat32(args[2])}
		v2 := rl.Vector2{X: toFloat32(args[3]), Y: toFloat32(args[4])}
		v3 := rl.Vector2{X: toFloat32(args[5]), Y: toFloat32(args[6])}
		rl.ImageDrawTriangleEx(img, v1, v2, v3, argsToColor(args, 7), argsToColor(args, 11), argsToColor(args, 15))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawTriangleLines", func(args []interface{}) (interface{}, error) {
		if len(args) < 10 {
			return nil, fmt.Errorf("ImageDrawTriangleLines requires (imageId, x1,y1, x2,y2, x3,y3, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		v1 := rl.Vector2{X: toFloat32(args[1]), Y: toFloat32(args[2])}
		v2 := rl.Vector2{X: toFloat32(args[3]), Y: toFloat32(args[4])}
		v3 := rl.Vector2{X: toFloat32(args[5]), Y: toFloat32(args[6])}
		rl.ImageDrawTriangleLines(img, v1, v2, v3, argsToColor(args, 7))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawTriangleFan", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ImageDrawTriangleFan requires (imageId, pointCount, x1,y1, x2,y2, ..., r,g,b,a)")
		}
		pointCount := int(toInt32(args[1]))
		if pointCount <= 0 || len(args) < 2+pointCount*2+4 {
			return nil, fmt.Errorf("ImageDrawTriangleFan needs pointCount and pointCount*2 coords plus color")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		points := make([]rl.Vector2, pointCount)
		for i := 0; i < pointCount; i++ {
			points[i] = rl.Vector2{X: toFloat32(args[2+i*2]), Y: toFloat32(args[2+i*2+1])}
		}
		colorOffset := 2 + pointCount*2
		rl.ImageDrawTriangleFan(img, points, argsToColor(args, colorOffset))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawTriangleStrip", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("ImageDrawTriangleStrip requires (imageId, pointCount, x1,y1, ..., r,g,b,a)")
		}
		pointCount := int(toInt32(args[1]))
		if pointCount <= 0 || len(args) < 2+pointCount*2+4 {
			return nil, fmt.Errorf("ImageDrawTriangleStrip needs pointCount and pointCount*2 coords plus color")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		points := make([]rl.Vector2, pointCount)
		for i := 0; i < pointCount; i++ {
			points[i] = rl.Vector2{X: toFloat32(args[2+i*2]), Y: toFloat32(args[2+i*2+1])}
		}
		colorOffset := 2 + pointCount*2
		rl.ImageDrawTriangleStrip(img, points, argsToColor(args, colorOffset))
		return nil, nil
	})
	v.RegisterForeign("ImageDraw", func(args []interface{}) (interface{}, error) {
		if len(args) < 13 {
			return nil, fmt.Errorf("ImageDraw requires (dstId, srcId, srcX,srcY,srcW,srcH, dstX,dstY,dstW,dstH, tintR,g,b,a)")
		}
		imageMu.Lock()
		dst, ok1 := images[toString(args[0])]
		src, ok2 := images[toString(args[1])]
		imageMu.Unlock()
		if !ok1 || dst == nil {
			return nil, fmt.Errorf("unknown dst image id: %s", toString(args[0]))
		}
		if !ok2 || src == nil {
			return nil, fmt.Errorf("unknown src image id: %s", toString(args[1]))
		}
		srcRec := rl.Rectangle{X: toFloat32(args[2]), Y: toFloat32(args[3]), Width: toFloat32(args[4]), Height: toFloat32(args[5])}
		dstRec := rl.Rectangle{X: toFloat32(args[6]), Y: toFloat32(args[7]), Width: toFloat32(args[8]), Height: toFloat32(args[9])}
		rl.ImageDraw(dst, src, srcRec, dstRec, argsToColor(args, 10))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawText", func(args []interface{}) (interface{}, error) {
		if len(args) < 7 {
			return nil, fmt.Errorf("ImageDrawText requires (imageId, posX, posY, text, fontSize, r, g, b, a)")
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		// raylib: ImageDrawText(dst *Image, posX, posY int32, text string, fontSize int32, col color.RGBA)
		rl.ImageDrawText(img, toInt32(args[1]), toInt32(args[2]), toString(args[3]), toInt32(args[4]), argsToColor(args, 5))
		return nil, nil
	})
	v.RegisterForeign("ImageDrawTextEx", func(args []interface{}) (interface{}, error) {
		if len(args) < 9 {
			return nil, fmt.Errorf("ImageDrawTextEx requires (imageId, fontId, posX, posY, text, fontSize, spacing, r, g, b, a)")
		}
		fontMu.Lock()
		font, okF := fonts[toString(args[1])]
		fontMu.Unlock()
		if !okF {
			return nil, fmt.Errorf("unknown font id: %s", toString(args[1]))
		}
		imageMu.Lock()
		img, ok := images[toString(args[0])]
		imageMu.Unlock()
		if !ok || img == nil {
			return nil, fmt.Errorf("unknown image id: %s", toString(args[0]))
		}
		// ImageDrawTextEx(dst *Image, position Vector2, font Font, text string, fontSize, spacing float32, col color.RGBA)
		pos := rl.Vector2{X: toFloat32(args[2]), Y: toFloat32(args[3])}
		rl.ImageDrawTextEx(img, pos, font, toString(args[4]), toFloat32(args[5]), toFloat32(args[6]), argsToColor(args, 7))
		return nil, nil
	})
}
