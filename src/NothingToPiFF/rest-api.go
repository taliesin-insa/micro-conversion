package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func generatePiFF(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	result := convertListToPiFF(string(reqBody))

	// get struct version
	//var piFFResult PiFFList
	//err = json.Unmarshal(result, &piFFResult)
	//checkError(err)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(result)
	checkError(err)
}

func main() {
	// Define the routing
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/convert", generatePiFF).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
