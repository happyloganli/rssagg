package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {

	if code > 499 {
		log.Println("Responding with error: " + message)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{
		Error: message,
	})

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	data, err := json.Marshal(payload) // Marshal payload to a json string and return a byte[]

	if err != nil {
		log.Printf("Failed to marshal JSON response: %s\n", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(data)
	if err != nil {
		return
	}
}
