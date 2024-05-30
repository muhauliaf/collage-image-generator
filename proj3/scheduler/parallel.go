package scheduler

import (
	"fmt"
	"image"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"proj3/png"
	"strings"
	"time"
)

// tileGenerator generates tile image from directory entry
func tileGenerator(config *Config, fileChannel <-chan fs.DirEntry, tileChannel chan<- *png.Image) {
	for {
		file, more := <-fileChannel
		if !more {
			break
		}
		filename := file.Name()
		ext := strings.ToLower(filepath.Ext(filename))
		if ext != ".png" {
			tileChannel <- nil
			continue
		}
		tilePath := filepath.Join(config.TilesDir, filename)
		tileImg, err := png.Load(tilePath)
		if err != nil {
			tileChannel <- nil
			continue
		}
		tileImg = tileImg.Resize(config.TileSize, config.TileSize)
		tileChannel <- tileImg
	}
}

// mosaicWorker applies color effects to input image in a specific tile position
func mosaicWorker(config *Config, outImg *png.Image, tileImgs []*png.Image, rectChannel <-chan *image.Rectangle, boolChannel chan<- bool) {
	for {
		bounds, more := <-rectChannel
		if !more {
			break
		}
		refImg := outImg.Subsize(*bounds)
		tileImg := tileImgs[rand.Intn(len(tileImgs))]
		colorTileImg := tileImg.ColorTransfer(refImg)
		for x := 0; x < bounds.Dx(); x++ {
			for y := 0; y < bounds.Dy(); y++ {
				blendTileColor := png.ColorBlend(
					png.ColortoRGBA64(tileImg.At(x, y)),
					png.ColortoRGBA64(colorTileImg.At(x, y)),
					config.Blendin,
				)
				outColor := png.ColorBlend(
					png.ColortoRGBA64(outImg.At(x+bounds.Min.X, y+bounds.Min.Y)),
					blendTileColor,
					config.Intensity,
				)
				outImg.Set(x+bounds.Min.X, y+bounds.Min.Y, outColor)
			}
		}
		boolChannel <- true
	}
}

// RunParallel runs the sequential version of the mosaic collage generator using channel
func RunParallel(config *Config) {
	// First part: creating upscaled input image and resized tile images
	var startTime time.Time
	var endTime float64
	startTime = time.Now()

	inImg, err := png.Load(config.InImg)
	ErrorCheck(err)
	files, err := os.ReadDir(config.TilesDir)
	ErrorCheck(err)

	outImg := inImg.Resize(inImg.Bounds().Dx()*config.Upscale, inImg.Bounds().Dy()*config.Upscale)
	bounds := outImg.Bounds()

	tileImgs := []*png.Image{}

	fileChannel := make(chan fs.DirEntry, len(files))
	tileChannel := make(chan *png.Image, config.Threads)

	// pushes tile generating tasks to a channel
	for _, file := range files {
		fileChannel <- file
	}
	close(fileChannel)

	// runs the tile generators
	for i := 0; i < config.Threads; i++ {
		go tileGenerator(config, fileChannel, tileChannel)
	}

	// extracting the result into an array
	for i := 0; i < len(files); i++ {
		tileImg := <-tileChannel
		if tileImg != nil {
			tileImgs = append(tileImgs, tileImg)
		}
	}
	close(tileChannel)

	endTime = time.Since(startTime).Seconds()
	fmt.Printf("%.2f\n", endTime)

	// Second part: applying color transfer to tile images, then add it input image position
	startTime = time.Now()

	rects := []*image.Rectangle{}
	for x0 := bounds.Min.X; x0 < bounds.Max.X; x0 += config.TileSize {
		for y0 := bounds.Min.Y; y0 < bounds.Max.Y; y0 += config.TileSize {
			x1 := min(x0+config.TileSize, bounds.Max.X)
			y1 := min(y0+config.TileSize, bounds.Max.Y)
			rect := image.Rect(x0, y0, x1, y1)
			rects = append(rects, &rect)
		}
	}

	rectChannel := make(chan *image.Rectangle, len(rects))
	boolChannel := make(chan bool, config.Threads)

	// pushes tile positions to channel
	for _, rect := range rects {
		rectChannel <- rect
	}
	close(rectChannel)

	// runs the mosaic worker
	for i := 0; i < config.Threads; i++ {
		go mosaicWorker(config, outImg, tileImgs, rectChannel, boolChannel)
	}

	// waiting until all tasks are finished
	for i := 0; i < len(rects); i++ {
		<-boolChannel
	}
	close(boolChannel)

	outImg.Save(config.OutImg)

	endTime = time.Since(startTime).Seconds()
	fmt.Printf("%.2f\n", endTime)
}
