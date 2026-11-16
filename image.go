package ccsds_spike

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Image struct {
	nZ          int       `toml:"n_z"`
	nX          int       `toml:"n_x"`
	nY          int       `toml:"n_y"`
	isSignedInt bool      `toml:"is_signed_int"`
	image_bytes [][][]int `toml:"image_bytes"`
}

func (ccsds *CCSDS) sampleAt(z, y, x int) int {
	return ccsds.image.image_bytes[z][y][x]
}

func (ccsds *CCSDS) fetchImage(path string) (Image, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return Image{}, fmt.Errorf("Error while reading image!", err)
	}
	var image Image
	err = toml.Unmarshal(fileBytes, &image)
	if err != nil {
		return Image{}, fmt.Errorf("Error while parsing image", err)
	}
	return image, nil
}
