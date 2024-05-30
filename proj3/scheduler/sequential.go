package scheduler

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"path/filepath"
	"proj3/png"
	"strings"
	"time"
)

// RunSequential runs the sequential version of the mosaic collage generator
func RunSequential(config *Config) {
	// First part: creating upscaled input image and resized tile images
	var startTime time.Time
	var endTime float64
	startTime = time.Now()

	// loads input files
	inImg, err := png.Load(config.InImg)
	ErrorCheck(err)
	files, err := os.ReadDir(config.TilesDir)
	ErrorCheck(err)

	// resizes input file
	outImg := inImg.Resize(inImg.Bounds().Dx()*config.Upscale, inImg.Bounds().Dy()*config.Upscale)
	bounds := outImg.Bounds()

	// loads tile images
	tileImgs := []*png.Image{}
	for _, file := range files {
		filename := file.Name()
		ext := strings.ToLower(filepath.Ext(filename))
		if ext != ".png" {
			continue
		}
		tilePath := filepath.Join(config.TilesDir, filename)
		tileImg, err := png.Load(tilePath)
		if err != nil {
			continue
		}
		tileImg = tileImg.Resize(config.TileSize, config.TileSize)
		tileImgs = append(tileImgs, tileImg)
	}
	endTime = time.Since(startTime).Seconds()
	fmt.Printf("%.2f\n", endTime)

	// Second part: applying color transfer to tile images, then add it input image position
	startTime = time.Now()

	// For each tile sized square in upscaled
	for x0 := bounds.Min.X; x0 < bounds.Max.X; x0 += config.TileSize {
		for y0 := bounds.Min.Y; y0 < bounds.Max.Y; y0 += config.TileSize {
			// selects a random image from tiles
			tileImg := tileImgs[rand.Intn(len(tileImgs))]
			x1 := min(x0+config.TileSize, bounds.Max.X)
			y1 := min(y0+config.TileSize, bounds.Max.Y)

			// extracts subimage at tile position
			refBounds := image.Rect(x0, y0, x1, y1)
			refImg := outImg.Subsize(refBounds)

			// applies color transfer to the tile image based on imput image
			colorTileImg := tileImg.ColorTransfer(refImg)

			// updates colored tile image to input image with weights
			for x := 0; x < x1-x0; x++ {
				for y := 0; y < y1-y0; y++ {
					blendTileColor := png.ColorBlend(
						png.ColortoRGBA64(tileImg.At(x, y)),
						png.ColortoRGBA64(colorTileImg.At(x, y)),
						config.Blendin,
					)
					outColor := png.ColorBlend(
						png.ColortoRGBA64(outImg.At(x+x0, y+y0)),
						blendTileColor,
						config.Intensity,
					)
					outImg.Set(x+x0, y+y0, outColor)
				}
			}
		}
	}

	// Saves output image
	outImg.Save(config.OutImg)
	endTime = time.Since(startTime).Seconds()
	fmt.Printf("%.2f\n", endTime)
}
