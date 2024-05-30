package png

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	*image.RGBA64
}

type Histogram struct {
}

func NewImage(width int, height int) *Image {
	bounds := image.Rect(0, 0, width, height)
	return &Image{image.NewRGBA64(bounds)}
}

func LoadFromImage(imgOrig image.Image) *Image {
	var img *Image
	bounds := imgOrig.Bounds()
	img = NewImage(bounds.Dx(), bounds.Dy())
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			r, g, b, a := imgOrig.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			img.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}
	return img
}

func LoadToImage(filePath string) (image.Image, error) {
	inReader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer inReader.Close()
	img, err := png.Decode(inReader)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func Load(filePath string) (*Image, error) {
	imgOrig, err := LoadToImage(filePath)
	if err != nil {
		return nil, err
	}
	img := LoadFromImage(imgOrig)
	return img, nil
}

func LoadDir(dirPath string) ([]*Image, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	imgs := []*Image{}
	for _, file := range files {
		filename := file.Name()
		ext := strings.ToLower(filepath.Ext(filename))
		if ext != ".png" {
			continue
		}
		imagePath := filepath.Join(dirPath, filename)
		img, err := Load(imagePath)
		if err != nil {
			continue
		}
		imgs = append(imgs, img)
	}
	return imgs, nil
}

func (img *Image) Clone(fill bool) *Image {
	bounds := img.Bounds()
	newImg := &Image{image.NewRGBA64(bounds)}
	if fill {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				newImg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
			}
		}
	}
	return newImg
}

func (img *Image) Save(filePath string) error {

	outWriter, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outWriter.Close()

	err = png.Encode(outWriter, img)
	if err != nil {
		return err
	}
	return nil
}
