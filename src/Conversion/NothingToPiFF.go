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
	Images []string
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

type ConversionError struct {
	URL         string
	Type        error
	Description string
}

type PiFFList struct {
	List   []PiFFStruct
	Errors []ConversionError
}

func convertListToPiFF(imagesPath []string) []byte {
	// initialization of returned variables
	var listOfPiFF []PiFFStruct
	var listOfErrors []ConversionError

	// create piFF for each image
	for i, im := range imagesPath {
		piFF, err := createPiFF(im)
		if err != nil {
			listOfErrors = append(listOfErrors, ConversionError{
				URL:         imagesPath[i],
				Type:        err,
				Description: err.Error(),
			})
		} else {
			listOfPiFF = append(listOfPiFF, piFF)
		}
	}

	// create the final Json to return
	piFFList := PiFFList{
		List:   listOfPiFF,
		Errors: listOfErrors,
	}
	result, err := json.MarshalIndent(piFFList, "", "     ")
	checkError(err)

	return result
}

// Fill struct, convert to json and write in file
func createPiFF(imagePath string) (PiFFStruct, error) {
	// get dimensions of image
	height, width, err := getDimensions(imagePath)
	if err != nil {
		return PiFFStruct{}, err
	}

	// fill PiFF struct
	PiFFData := PiFFStruct{
		Meta: Meta{
			Type: "line",
			URL:  imagePath,
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

	return PiFFData, nil
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
