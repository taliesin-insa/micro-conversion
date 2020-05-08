package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestGeneratePiFFNilBody(t *testing.T) {
	assert := assert.New(t)

	request, err := http.NewRequest("POST", "/convert/nothing", nil)
	if err != nil {
		t.Errorf("[TEST_ERROR] Create request: %v", err.Error())
	}

	recorder := httptest.NewRecorder()
	generatePiFF(recorder, request)

	assert.Equal(http.StatusBadRequest, recorder.Code)

	assert.Equal("[MICRO-CONVERSION] Request body is null", string(recorder.Body.Bytes()))
}

func TestGeneratePiFFWrongBody(t *testing.T) {
	assert := assert.New(t)

	request, err := http.NewRequest("POST", "/convert/nothing", bytes.NewBuffer([]byte("Wrong body format")))
	if err != nil {
		t.Errorf("[TEST_ERROR] Create request: %v", err.Error())
	}

	recorder := httptest.NewRecorder()
	generatePiFF(recorder, request)

	assert.Equal(http.StatusBadRequest, recorder.Code)

	assert.Equal("[MICRO-CONVERSION] Couldn't unmarshal body", string(recorder.Body.Bytes()))
}

func TestGeneratePiFFUnknownImage(t *testing.T) {
	assert := assert.New(t)

	unknownImage := RequestDataNothing{Path: "unknown/image.png"}
	unknownImageJSON, err := json.Marshal(unknownImage)
	if err != nil {
		t.Errorf("[TEST_ERROR] Marshal request body: %v", err.Error())
	}

	request, err := http.NewRequest("POST", "/convert/nothing", bytes.NewBuffer(unknownImageJSON))
	if err != nil {
		t.Errorf("[TEST_ERROR] Create request: %v", err.Error())
	}

	recorder := httptest.NewRecorder()
	generatePiFF(recorder, request)

	assert.Equal(http.StatusBadRequest, recorder.Code)

	assert.Equal("[MICRO-CONVERSION] Couldn't open image", string(recorder.Body.Bytes()))
}

func TestGeneratePiFF(t *testing.T) {
	assert := assert.New(t)

	// create an image for the test
	width := 200
	height := 100

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	image := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// set color for each pixel
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			image.Set(x, y, color.White)
		}
	}

	// save temporary the image
	imageFile, err := ioutil.TempFile("", "MICRO_CONVERSION_TEST")
	if err != nil {
		t.Errorf("[TEST_ERROR] Create image: %v", err.Error())
	}
	imagePath := imageFile.Name()

	err = png.Encode(imageFile, image)
	if err != nil {
		imageFile.Close()
		t.Errorf("[TEST_ERROR] Encode image: %v", err.Error())
	}

	if imageFile.Close() != nil {
		t.Errorf("[TEST_ERROR] Close image: %v", err.Error())
	}

	// tests with this image
	imageRequest := RequestDataNothing{Path: imagePath}
	imageRequestJSON, err := json.Marshal(imageRequest)
	if err != nil {
		t.Errorf("[TEST_ERROR] Marshal request body: %v", err.Error())
	}

	request, err := http.NewRequest("POST", "/convert/nothing", bytes.NewBuffer(imageRequestJSON))
	if err != nil {
		t.Errorf("[TEST_ERROR] Create request: %v", err.Error())
	}

	recorder := httptest.NewRecorder()
	generatePiFF(recorder, request)

	// status test
	assert.Equal(http.StatusOK, recorder.Code)

	// body format test
	var reqData PiFFStruct
	err = json.Unmarshal(recorder.Body.Bytes(), &reqData)
	assert.Nil(err, "Handler returned wrong format body: "+string(recorder.Body.Bytes()))

	// Location content test
	locations := reqData.Location
	assert.Equal(1, len(locations))

	dimensions := locations[0].Polygon
	dimensionsTest := [][2]int{
		{0, 0},
		{height, 0},
		{height, width},
		{0, width},
	}
	assert.True(reflect.DeepEqual(dimensionsTest, dimensions))

	// delete the useless image
	err = os.Remove(imagePath)
	if err != nil {
		t.Errorf("[TEST_ERROR] Delete image: %v", err.Error())
	}
}
