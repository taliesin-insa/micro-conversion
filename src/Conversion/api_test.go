package main

import (
	"bytes"
	"encoding/json"
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
	request, err := http.NewRequest("POST", "/convert/nothing", nil)
	if err != nil {
		t.Errorf("[TEST_ERROR] Create request: %v", err.Error())
	}

	recorder := httptest.NewRecorder()
	generatePiFF(recorder, request)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("[TEST_ERROR] Handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if message := string(recorder.Body.Bytes()); message != "[MICRO-CONVERSION] Request body is null" {
		t.Errorf("[TEST_ERROR] Handler returned wrong response body: got %v want %v",
			message, "[MICRO-CONVERSION] Request body is null")
	}

}

func TestGeneratePiFFWrongBody(t *testing.T) {
	request, err := http.NewRequest("POST", "/convert/nothing", bytes.NewBuffer([]byte("Wrong body format")))
	if err != nil {
		t.Errorf("[TEST_ERROR] Create request: %v", err.Error())
	}

	recorder := httptest.NewRecorder()
	generatePiFF(recorder, request)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("[TEST_ERROR] Handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if message := string(recorder.Body.Bytes()); message != "[MICRO-CONVERSION] Couldn't unmarshal body" {
		t.Errorf("[TEST_ERROR] Handler returned wrong response body: got %v want %v",
			message, "[MICRO-CONVERSION] Couldn't unmarshal body")
	}
}

func TestGeneratePiFFUnknownImage(t *testing.T) {
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

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("[TEST_ERROR] Handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if message := string(recorder.Body.Bytes()); message != "[MICRO-CONVERSION] Couldn't open image" {
		t.Errorf("[TEST_ERROR] Handler returned wrong response body: got %v want %v",
			message, "[MICRO-CONVERSION] Couldn't open image")
	}
}

func TestGeneratePiFF(t *testing.T) {
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
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("[TEST_ERROR] Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// body format test
	var reqData PiFFStruct
	err = json.Unmarshal(recorder.Body.Bytes(), &reqData)
	if err != nil {
		t.Errorf("[TEST_ERROR] Handler returned wrong format body (got %v): %v",
			string(recorder.Body.Bytes()), err.Error())
	}

	// Location content test
	locations := reqData.Location
	if len(locations) != 1 {
		t.Errorf("[TEST_ERROR] Handler returned wrong body content (length of Location): got %v want %v",
			len(locations), 1)
	}

	dimensions := locations[0].Polygon
	dimensionsTest := [][2]int{
		{0, 0},
		{height, 0},
		{height, width},
		{0, width},
	}
	if !reflect.DeepEqual(dimensionsTest, dimensions) {
		t.Errorf("[TEST_ERROR] Handler returned wrong body content (content of Polygon): got %v want %v",
			dimensions, dimensionsTest)
	}

	// delete the useless image
	err = os.Remove(imagePath)
	if err != nil {
		t.Errorf("[TEST_ERROR] Delete image: %v", err.Error())
	}
}
