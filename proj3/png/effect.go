package png

import (
	"image"
	"image/color"
	"math"
	"proj3/utils"
)

func (img *Image) Histogram() [3][256]float64 {
	bounds := (*img).Bounds()
	var hist [3][256]float64
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := (*img).At(x, y).RGBA()
			hist[0][r>>8]++
			hist[1][g>>8]++
			hist[2][b>>8]++
		}
	}
	totalPixels := float64(bounds.Dx() * bounds.Dy())

	for i := 0; i < 3; i++ {
		for j := 0; j < 256; j++ {
			hist[i][j] /= totalPixels
		}
	}

	return hist
}

func (img *Image) CDF() [3][256]float64 {
	bounds := (*img).Bounds()
	var histCDF [3][256]float64
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := (*img).At(x, y).RGBA()
			histCDF[0][r>>8]++
			histCDF[1][g>>8]++
			histCDF[2][b>>8]++
		}
	}
	totalPixels := float64(bounds.Dx() * bounds.Dy())

	for i := 0; i < 3; i++ {
		var sum float64
		for j := 0; j < 256; j++ {
			sum += histCDF[i][j]
			histCDF[i][j] = sum / totalPixels
		}
	}
	return histCDF
}

func (img *Image) MapPixels(pixels [3][256]float64) *Image {
	bounds := img.Bounds()
	newImg := NewImage(bounds.Dx(), bounds.Dy())
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			r, g, b, a := (*img).At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			r = uint32(math.Round(pixels[0][r>>8])) * 256
			g = uint32(math.Round(pixels[1][g>>8])) * 256
			b = uint32(math.Round(pixels[2][b>>8])) * 256
			newImg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}
	return newImg
}

func (img *Image) ColorTransfer(imgRef *Image) *Image {
	var matchImg *Image
	inCDF := img.CDF()
	refCDF := imgRef.CDF()
	var pixels [3][256]float64
	var newPixels [3][256]float64
	for i := 0; i < 3; i++ {
		pixels[i] = [256]float64(utils.ArangeFloat(256))
		newPixels[i] = [256]float64(utils.Interpolate(inCDF[i][:], refCDF[i][:], pixels[i][:]))
	}
	matchImg = img.MapPixels(newPixels)
	return matchImg
}

func ColortoRGBA64(clr color.Color) *color.RGBA64 {
	r, g, b, a := clr.RGBA()
	return &color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func ColorBlend(srcColor *color.RGBA64, dstColor *color.RGBA64, dstWeight float64) *color.RGBA64 {
	if dstWeight <= 0.0 {
		return srcColor
	} else if dstWeight >= 1.0 {
		return dstColor
	} else {
		r2 := math.Round((float64(srcColor.R) * (1 - dstWeight)) + (float64(dstColor.R) * dstWeight))
		g2 := math.Round((float64(srcColor.G) * (1 - dstWeight)) + (float64(dstColor.G) * dstWeight))
		b2 := math.Round((float64(srcColor.B) * (1 - dstWeight)) + (float64(dstColor.B) * dstWeight))
		a2 := math.Round((float64(srcColor.A) * (1 - dstWeight)) + (float64(dstColor.A) * dstWeight))
		return &color.RGBA64{uint16(r2), uint16(g2), uint16(b2), uint16(a2)}
	}
}

func (img *Image) Resize(width int, height int) *Image {
	newImg := NewImage(width, height)
	boundsOri := img.Bounds()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			x0 := int(float64(x) * float64(boundsOri.Dx()) / float64(width))
			y0 := int(float64(y) * float64(boundsOri.Dy()) / float64(height))
			r, g, b, a := img.At(x0, y0).RGBA()
			newImg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}
	return newImg
}

func (img *Image) Subsize(bounds image.Rectangle) *Image {
	newImg := NewImage(bounds.Dx(), bounds.Dy())
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			r, g, b, a := (*img).At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			newImg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}
	return newImg
}
