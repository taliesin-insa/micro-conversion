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
		log.Fatal(err)
		return
	}

	// converts body to json
	var reqData RequestDataNothing
	err = json.Unmarshal(reqBody, &reqData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Wrong request body format: this microservice needs a JSON which attribute 'Images' is an array of string"))
		return
	}

	// get PiFF list
	result := convertListToPiFF(reqData.Images)

	// send answer
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(result)
	checkError(err)
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
