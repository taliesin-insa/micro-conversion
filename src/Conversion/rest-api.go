package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
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

func generatePiFF(w http.ResponseWriter, r *http.Request) {
	// get request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Read body: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[MICRO-CONVERSION] Couldn't read body"))
		return
	}

	log.Printf("[BODY] %v", reqBody)

	// converts body to json
	var reqData RequestDataNothing
	err = json.Unmarshal(reqBody, &reqData)
	if err != nil {
		log.Printf("[ERROR] Unmarshal body: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[MICRO-CONVERSION] Couldn't unmarshal body"))
		return
	}

	// get dimensions of image
	file, err := os.Open(reqData.Path)
	if err != nil {
		log.Printf("[ERROR] Open image: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[MICRO-CONVERSION] Couldn't open image"))
		return
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Printf("[ERROR] Read image: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[MICRO-CONVERSION] Couldn't read image"))
		return
	}

	width := image.Width
	height := image.Height

	err = file.Close()
	if err != nil {
		log.Printf("[ERROR] Close image: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[MICRO-CONVERSION] Couldn't close image"))
		return
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
	if err != nil {
		log.Printf("[ERROR] Marshal piFF: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[MICRO-CONVERSION] Couldn't marshal piFF"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// function to test whether docker file is correctly built
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[MICRO-CONVERSION] Welcome home!")
}

func main() {
	// Define the routing
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/convert/nothing", generatePiFF).Methods("POST")
	router.HandleFunc("/convert", homeLink).Methods("GET")
	log.Fatal(http.ListenAndServe(":12345", router))
}
