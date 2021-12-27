package imagecompression

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"os"

	"github.com/nfnt/resize"
)

func ImageCompress(
	getReadSizeFile func() (io.Reader, error),
	getDecodeFile func() (*os.File, error),
	to string,
	quality,
	base int,
	format string) (err error) {
	// read file
	fileOrigin, err := getDecodeFile()
	if err != nil {
		return
	}

	defer func() {
		if fileOrigin != nil {
			fileOrigin.Close()
		}
	}()

	var (
		origin    image.Image
		config    image.Config
		temp      io.Reader
		typeImage int64
	)

	// read file size
	_, err = getReadSizeFile()
	if err != nil {
		return
	}

	format = strings.ToLower(format)
	switch format {
	case "jpg":
		fallthrough
	case "jpeg":
		{
			typeImage = 1
			origin, err = jpeg.Decode(fileOrigin)
			if err != nil {
				return
			}

			temp, err = getReadSizeFile()
			if err != nil {
				return
			}

			config, err = jpeg.DecodeConfig(temp)
			if err != nil {
				return
			}
		}

	case "png":
		{
			typeImage = 0
			origin, err = png.Decode(fileOrigin)
			if err != nil {
				return
			}
			temp, err = getReadSizeFile()
			if err != nil {
				return
			}
			config, err = png.DecodeConfig(temp)
			if err != nil {
				return
			}
		}

	default:
		err = fmt.Errorf("image format %s not supportted", format)
		return
	}

	// resize the image
	width := uint(base)
	height := uint(base * config.Height / config.Width)
	canvas := resize.Thumbnail(width, height, origin, resize.Lanczos3)

	fileOut, err := os.Create(to)
	if err != nil {
		return
	}
	defer func() {
		if fileOut != nil {
			fileOut.Close()
		}
	}()

	switch typeImage {
	case 0:
		{
			err = png.Encode(fileOut, canvas)
			if err != nil {
				return
			}
		}
	case 1:
		{
			err = jpeg.Encode(fileOut, canvas, &jpeg.Options{Quality: quality})
			if err != nil {
				return
			}
		}
	default:
		err = fmt.Errorf("typeImage %d not supportted", typeImage)
		return
	}

	return
}
