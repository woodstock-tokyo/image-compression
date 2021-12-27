package imagecompression

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// TestImageCompression test image compression
func TestImageCompression(t *testing.T) {
	filePath := "akb.jpg"
	file, err := os.Open(filePath)
	if err != nil {
		t.Error(fmt.Errorf("failed to read file"))
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size() / 1024 // size in kilobytes
	fmt.Printf("original file size is %d kB\n", size)

	getReadSizeFileFn := func() (io.Reader, error) {
		return os.Open(filePath)
	}

	getDecodeFileFn := func() (*os.File, error) {
		return os.Open(filePath)
	}

	_, format, shortFileName := parseImageFormat(filePath)
	output := shortFileName + "_compress." + format
	quality := 80
	width := 200

	err = ImageCompress(getReadSizeFileFn, getDecodeFileFn, output, quality, width, format)
	if err != nil {
		t.Error(err)
	}

	outputPath := "akb_compress.jpg"
	outputFile, err := os.Open(outputPath)
	if err != nil {
		t.Error(fmt.Errorf("failed to read file"))
	}
	defer outputFile.Close()

	fileInfo, _ = outputFile.Stat()
	size = fileInfo.Size() / 1024 // size in kilobytes
	fmt.Printf("compressed file size is %d kB\n", size)
}

func parseImageFormat(path string) (filePath string, format string, shortFileName string) {
	temp := strings.Split(path, ".")
	if len(temp) <= 1 {
		return "", "", ""
	}
	mapRule := make(map[string]bool)
	mapRule["jpg"] = true
	mapRule["png"] = true
	mapRule["jpeg"] = true
	/** add more format here */

	if mapRule[temp[1]] {
		return path, temp[1], temp[0]
	} else {
		return "", "", ""
	}
}
