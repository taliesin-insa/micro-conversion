package main

import (
	"encoding/json"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"
)

type Meta struct {
	Type string
	URL  string
	Id   int
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
	Location Location
	Data     []Data
	Children []int
	Parent   int
}

type ConversionError struct {
	URL         string
	Description string
}

type PiFFList struct {
	List   []PiFFStruct
	Errors []ConversionError
}

//OpenFile fail
//readFile fail

//func main() {
//	if len(os.Args) != 2 {
//		fmt.Println("Usage: go run NothingToPiFF.go filePath")
//		os.Exit(1)
//	}
//
//	resultJson := convertListToPiFF(os.Args[1])
//	var piFFList PiFFList
//	err := json.Unmarshal(resultJson, &piFFList)
//	checkError(err)
//
//	for _, piFFFile := range piFFList.List {
//		fmt.Println(string(piFFFile))
//	}
//}

func convertListToPiFF(filePath string) []byte {
	// initialization of returned variables
	var listOfPiFF []PiFFStruct
	var listOfErrors []ConversionError

	// open the file
	dataByte, err := ioutil.ReadFile(filePath)
	if err != nil {
		listOfErrors = append(listOfErrors, ConversionError{
			URL:         filePath,
			Description: err.Error(),
		})
	}

	// get the list of images to convert
	data := string(dataByte)
	nameOfImages := strings.Split(data, "\n")
	nameOfImages = nameOfImages[:len(nameOfImages)-1]

	// create piFF for each image
	for _, im := range nameOfImages {
		piFF, err := createPiFF(im)
		if err != nil {
			listOfErrors = append(listOfErrors, ConversionError{
				URL:         filePath,
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
	// get name of image
	segments := strings.Split(imagePath, "/")
	imageName := segments[len(segments)-1]

	// get dimensions of image
	height, width, err := getDimensions(imagePath)
	if err != nil {
		return PiFFStruct{}, err
	}

	// fill PiFF struct
	PiFFData := PiFFStruct{
		Meta: Meta{
			Type: "line",
			URL:  imageName,
			Id:   0,
		},
		Location: Location{
			Type: "line",
			Polygon: [][2]int{
				{0, 0},
				{height, 0},
				{height, width},
				{0, width},
			},
			Id: "loc_0",
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
