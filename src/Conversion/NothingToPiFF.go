package main

import (
	"encoding/json"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type RequestDataNothing struct {
	Path string
}

type Meta struct {
	Type string
	URL  string
}

type Location struct {
	Type    string
	Polygon [][2]int
	Id      string
}

type Data struct {
	Type       string
	LocationId string
	Value      string
	Id         string
}

type PiFFStruct struct {
	Meta     Meta
	Location []Location
	Data     []Data
	Children []int
	Parent   int
}

// Fill struct, convert to json and write in file
func convertFromNothingToPiFF(imagePath string) ([]byte, error) {
	// get dimensions of image
	height, width, err := getDimensions(imagePath)
	if err != nil {
		return nil, err
	}

	// fill PiFF struct
	PiFFData := PiFFStruct{
		Meta: Meta{
			Type: "line",
			URL:  "",
		},
		Location: []Location{
			{Type: "line",
				Polygon: [][2]int{
					{0, 0},
					{height, 0},
					{height, width},
					{0, width},
				},
				Id: "loc_0",
			},
		},
		Data: []Data{
			{
				Type:       "line",
				LocationId: "loc_0",
				Value:      "",
				Id:         "0",
			},
		},
	}

	result, err := json.MarshalIndent(PiFFData, "", "     ")
	checkError(err)
	return result, nil
}

// Extract width and height of the image
func getDimensions(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}

	image, _, err := image.DecodeConfig(file)
	checkError(err)
	width := image.Width
	height := image.Height

	checkError(file.Close())

	return height, width, nil
}

// Check for unhandled errors
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
