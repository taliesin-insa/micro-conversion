package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeneratePiFFNilBody(t *testing.T) {
	request, err := http.NewRequest("POST", "/convert/nothing", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(generatePiFF)
	handler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if message := string(recorder.Body.Bytes()); message != "[MICRO-CONVERSION] Request body is null" {
		t.Errorf("handler returned wrong response body: got %v want %v",
			message, "[MICRO-CONVERSION] Request body is null")
	}

}

func TestGeneratePiFFWrongBody(t *testing.T) {
	request, err := http.NewRequest("POST", "/convert/nothing", bytes.NewBuffer([]byte("Wrong body format")))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(generatePiFF)
	handler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if message := string(recorder.Body.Bytes()); message != "[MICRO-CONVERSION] Couldn't unmarshal body" {
		t.Errorf("handler returned wrong response body: got %v want %v",
			message, "[MICRO-CONVERSION] Couldn't unmarshal body")
	}
}

func TestGeneratePiFFUnknownImage(t *testing.T) {
	unknownImage := RequestDataNothing{Path: "unknown/image.png"}
	unknownImageJSON, err := json.Marshal(unknownImage)
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest("POST", "/convert/nothing", bytes.NewBuffer(unknownImageJSON))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(generatePiFF)
	handler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if message := string(recorder.Body.Bytes()); message != "[MICRO-CONVERSION] Couldn't open image" {
		t.Errorf("handler returned wrong response body: got %v want %v",
			message, "[MICRO-CONVERSION] Couldn't open image")
	}
}

func TestGeneratePiFF(t *testing.T) {

}
