package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func generatePiFF(w http.ResponseWriter, r *http.Request) {
	// get request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Fatal(err.Error())
		return
	}

	// converts body to json
	var reqData RequestDataNothing
	err = json.Unmarshal(reqBody, &reqData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong request body format: this microservice needs a JSON which attribute 'Path' is a string"))
		log.Fatal(err.Error())
		return
	}

	// get PiFF list
	result, err := convertFromNothingToPiFF(reqData.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Fatal(err.Error())
		return

	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

// function to test whether docker file is correctly built
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	// Define the routing
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/convert/nothing", generatePiFF).Methods("POST")
	router.HandleFunc("/", homeLink).Methods("GET")
	log.Fatal(http.ListenAndServe(":12345", router))
}
