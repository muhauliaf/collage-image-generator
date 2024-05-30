package scheduler

import (
	"fmt"
	"os"
)

type Config struct {
	InImg     string
	OutImg    string
	TilesDir  string
	TileSize  int
	RunMode   string
	Threads   int
	Upscale   int
	Intensity float64
	Blendin   float64
}

// ErrorCheck checks for error, then if one exists, prints it then exit the application
func ErrorCheck(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

// Schedule runs the correct version based on the Mode field of the configuration value
func Schedule(config *Config) {
	if config.RunMode == "s" {
		RunSequential(config)
	} else if config.RunMode == "p" {
		RunParallel(config)
	} else if config.RunMode == "w" {
		RunWorkSteal(config)
	}
}
