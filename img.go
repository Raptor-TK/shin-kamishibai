package main

import (
	"bufio"
	"bytes"
	"image"
	"image/jpeg"
	"io"
)

// ImageThumb create thumbnail while maintaining ratio
func ImageThumb(reader io.Reader) ([]byte, error) {
	// clone so can use again
	var b bytes.Buffer
	reader2 := io.TeeReader(reader, &b)
	reader3 := bufio.NewReader(&b)

	maxW, maxH := float64(320), float64(320) // maximum allowed thumbnail dimension
	var imgW, imgH float64                   // original image width, height
	var thmW, thmH int                       // thumbnail width, height
	var ratio float64                        // image w/h ratio

	// get image dimension
	m, _, err := image.Decode(reader2)
	if err != nil {
		return nil, err
	}
	bounds := m.Bounds()
	imgW, imgH = float64(bounds.Max.X), float64(bounds.Max.Y)
	// fmt.Println(imgW, imgH)
	// image ratio
	ratio = float64(imgW) / float64(imgH)

	if maxW >= maxH {
		if imgW > imgH {
			thmW = int(maxW)
			thmH = int(float64(thmW) / ratio)
		} else {
			thmH = int(maxH)
			thmW = int(float64(thmH) * ratio)
		}
	} else {
		if imgW < imgH {
			thmW = int(maxW)
			thmH = int(float64(thmW) * ratio)
		} else {
			thmH = int(maxH)
			thmW = int(float64(thmH) / ratio)
		}
	}

	return ImageResize(reader3, thmW, thmH)
}

// ImageResize resize image to specific width, height
func ImageResize(reader io.Reader, owidth int, oheight int) ([]byte, error) {
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	bounds := m.Bounds()

	// fmt.Printf("%+v", bounds)
	// fmt.Println("dimension", bounds.Max.X, "x", bounds.Max.Y)

	// ratio
	var rx, ry float32
	rx = float32(owidth) / float32(bounds.Max.X)
	ry = float32(oheight) / float32(bounds.Max.Y)
	// fmt.Println("ratio x, y", rx, ry)

	newImg := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{owidth, oheight},
		},
	)

	for x := 0; x < owidth; x++ {
		for y := 0; y < oheight; y++ {
			// imported image cord
			ix := int(float32(x) / rx)
			iy := int(float32(y) / ry)

			// r, g, b, a := m.At(x, y).RGBA()
			// rgba := m.At(x, y)
			rgba := m.At(ix, iy)

			// fmt.Println("xy:", x, ix, y, iy, rgba)

			newImg.Set(x, y, rgba)
			// newImg.Set(x, y, color.White)
		}
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	opts := &jpeg.Options{Quality: 50}
	jpeg.Encode(writer, newImg, opts)

	return b.Bytes(), nil
}
