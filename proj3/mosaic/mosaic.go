package main

import (
	"flag"
	"fmt"
	"os"
	"proj3/scheduler"
)

// ErrorExit prints error and usage then exit the application
func ErrorExit(errMessage string) {
	fmt.Println("ERROR:", errMessage)
	fmt.Println("FLAGS:")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	inImg := flag.String("i", "", "Path to the input image")
	outImg := flag.String("o", "", "Path to the output image")
	tilesDir := flag.String("d", "", "Path to the mosaic tiles directory")
	tileSize := flag.Int("s", 0, "Size of mosaic tiles in pixels. Must be positive")
	upscale := flag.Int("U", 1, "Input image upscaling in integer. Must be positive")
	intensity := flag.Float64("I", 0.8, "intensity of mosaic images in float (0.0 - 1.0). 1.0 for full mosaic images, 0.0 for input image only")
	blendin := flag.Float64("B", 0.8, "intensity of tile images color blend-in in float (0.0 - 1.0). 1.0 for full blend in, 0.0 for tiles image only")
	runMode := flag.String("M", "s", "running mode: s=sequential(default), p=parallel, w=parallel with work steal")
	threads := flag.Int("T", 1, "Number of goroutines. ignored if sequential. Must be positive")

	flag.Parse()

	if *inImg == "" || *outImg == "" || *tilesDir == "" {
		ErrorExit("'i','o','d' flag is required")
	}
	if *tileSize < 1 {
		ErrorExit("'s' flag is required and must be positive")
	}
	if *upscale < 1 {
		ErrorExit("'U' must be positive")
	}
	if *intensity < 0.0 || *intensity > 1.0 {
		ErrorExit("'I' must be from 0.0 to 1.0")
	}
	if *blendin < 0.0 || *blendin > 1.0 {
		ErrorExit("'B' must be from 0.0 to 1.0")
	}
	if *runMode != "s" && *runMode != "p" && *runMode != "w" {
		ErrorExit("'M' must be: s, p, w")
	}
	if *threads < 1 {
		ErrorExit("'T' must be positive")
	}

	var config scheduler.Config = scheduler.Config{}
	config.InImg = *inImg
	config.OutImg = *outImg
	config.TilesDir = *tilesDir
	config.TileSize = *tileSize
	config.RunMode = *runMode
	config.Threads = *threads
	config.Upscale = *upscale
	config.Intensity = *intensity
	config.Blendin = *blendin
	scheduler.Schedule(&config)
}
