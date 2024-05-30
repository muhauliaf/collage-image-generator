package scheduler

import (
	"fmt"
	"image"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"proj3/deque"
	"proj3/png"
	"runtime"
	"strings"
	"time"
)

// generateTile generates tile image from directory entry
func generateTile(config *Config, file fs.DirEntry) *png.Image {
	filename := file.Name()
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".png" {
		return nil
	}
	tilePath := filepath.Join(config.TilesDir, filename)
	tileImg, err := png.Load(tilePath)
	if err != nil {
		return nil
	}
	tileImg = tileImg.Resize(config.TileSize, config.TileSize)
	return tileImg
}

// createMosaic applies color effects to input image in a specific tile position
func createMosaic(config *Config, bounds *image.Rectangle, outImg *png.Image, tileImgs []*png.Image) bool {
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
	return true
}

// workStealTileGenerator pops tasks from its deque, then tries to steals tasks from other deques if empty
func workStealTileGenerator(config *Config, id int, deques []*deque.BoundDeque, tileChannel chan<- *png.Image, done *bool) {
	task := deques[id].PopBottom()
	for {
		for task != nil {
			tileImg := generateTile(config, task.(fs.DirEntry))
			tileChannel <- tileImg
			task = deques[id].PopBottom()
		}
		for task == nil {
			if *done {
				return
			}
			runtime.Gosched()
			victim := rand.Intn(len(deques))
			if !deques[victim].IsEmpty() {
				task = deques[victim].PopTop()
			}
		}
	}
}

// workStealMosaicWorker pops tasks from its deque, then tries to steals tasks from other deques if empty
func workStealMosaicWorker(config *Config, id int, deques []*deque.BoundDeque, outImg *png.Image, tileImgs []*png.Image, boolChannel chan<- bool, done *bool) {
	task := deques[id].PopBottom()
	for {
		for task != nil {
			resp := createMosaic(config, task.(*image.Rectangle), outImg, tileImgs)
			boolChannel <- resp
			task = deques[id].PopBottom()
		}
		for task == nil {
			if *done {
				return
			}
			runtime.Gosched()
			victim := rand.Intn(len(deques))
			if !deques[victim].IsEmpty() {
				task = deques[victim].PopTop()
			}
		}
	}
}

// RunWorkSteal runs the fork join with work stealing version of the mosaic collage generator
func RunWorkSteal(config *Config) {
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

	tileDone := false
	tileImgs := []*png.Image{}
	tileChannel := make(chan *png.Image, len(files))

	// pushes tile generating tasks to deque in each thread
	deques := make([]*deque.BoundDeque, config.Threads)
	for i := 0; i < config.Threads; i++ {
		deques[i] = deque.NewBoundDeque((len(files) / config.Threads) + 1)
	}
	for i, file := range files {
		deques[i%config.Threads].PushBottom(file)
	}

	// runs the tile generators
	for i := 0; i < config.Threads; i++ {
		go workStealTileGenerator(config, i, deques, tileChannel, &tileDone)
	}

	for i := 0; i < len(files); i++ {
		tileImg := <-tileChannel
		if tileImg != nil {
			tileImgs = append(tileImgs, tileImg)
		}
	}

	tileDone = true
	close(tileChannel)

	endTime = time.Since(startTime).Seconds()
	fmt.Printf("%.2f\n", endTime)

	// Second part: applying color transfer to tile images, then add it input image position
	startTime = time.Now()

	rectDone := false
	rects := []*image.Rectangle{}
	boolChannel := make(chan bool, config.Threads)

	// pushes tile positions to deque in each thread
	for x0 := bounds.Min.X; x0 < bounds.Max.X; x0 += config.TileSize {
		for y0 := bounds.Min.Y; y0 < bounds.Max.Y; y0 += config.TileSize {
			x1 := min(x0+config.TileSize, bounds.Max.X)
			y1 := min(y0+config.TileSize, bounds.Max.Y)
			rect := image.Rect(x0, y0, x1, y1)
			rects = append(rects, &rect)
		}
	}

	deques = make([]*deque.BoundDeque, config.Threads)
	for i := 0; i < config.Threads; i++ {
		deques[i] = deque.NewBoundDeque((len(rects) / config.Threads) + 1)
	}
	for i, rect := range rects {
		deques[i%config.Threads].PushBottom(rect)
	}

	// runs the mosaic worker
	for i := 0; i < config.Threads; i++ {
		go workStealMosaicWorker(config, i, deques, outImg, tileImgs, boolChannel, &rectDone)
	}

	for i := 0; i < len(rects); i++ {
		<-boolChannel
	}

	rectDone = true
	close(boolChannel)

	outImg.Save(config.OutImg)

	endTime = time.Since(startTime).Seconds()
	fmt.Printf("%.2f\n", endTime)
}
