package main

import (
	"bufio"
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	_ "golang.org/x/image/webp" // Import the WebP package
	"io"
)

// ImageThumb create thumbnail image
func ImageThumb(reader io.Reader) ([]byte, error) {
	return ImageScale(reader, 320, 320)
}

// ImageScale resize image while maintaining ratio
func ImageScale(reader io.Reader, miW, miH int) ([]byte, error) {
	// clone so can use again
	var b bytes.Buffer
	reader2 := io.TeeReader(reader, &b)
	reader3 := bufio.NewReader(&b)

	maxW, maxH := float64(miW), float64(miH) // maximum allowed thumbnail dimension
	var imgW, imgH float64                   // original image width, height
	var thmW, thmH int                       // thumbnail width, height
	var ratio float64                        // image w/h ratio

	// get image dimension
	m, format, err := image.Decode(reader2)
	if err != nil {
		return nil, err
	}
	bounds := m.Bounds()
	imgW, imgH = float64(bounds.Max.X), float64(bounds.Max.Y)

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

	// If the original format is GIF, return the original data
	if format == "gif" {
		// Return the original GIF data
		return io.ReadAll(reader)
	}

	return ImageResize(reader3, thmW, thmH, format)
}

// ImageResize resize image to specific width, height
func ImageResize(reader io.Reader, owidth int, oheight int, format string) ([]byte, error) {
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	bounds := m.Bounds()

	// ratio
	var rx, ry float32
	rx = float32(owidth) / float32(bounds.Max.X)
	ry = float32(oheight) / float32(bounds.Max.Y)

	// new blank canvas
	newImg := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{owidth, oheight},
		},
	)

	// fill canvas using skip pixel (nearest neighbour)
	for x := 0; x < owidth; x++ {
		for y := 0; y < oheight; y++ {
			// imported image cord
			ix := int(float32(x) / rx)
			iy := int(float32(y) / ry)

			rgba := m.At(ix, iy)

			newImg.Set(x, y, rgba)
		}
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	// Encode jpeg
		opts := &jpeg.Options{Quality: 50}
		err = jpeg.Encode(writer, newImg, opts)

	if err != nil {
		return nil, err
	}

	writer.Flush() // Ensure all data is written to the buffer

	return b.Bytes(), nil
}
