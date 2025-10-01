package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerChirpValidation(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnError struct {
		Error string `json:"error"`
	}

	type returnVal struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		respBody := returnError{
			Error: "Chirp is too long",
		}
		respo, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling error: %s", err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(respo)
		return
	}

	respBody := returnVal{
		Valid: true,
	}
	respo, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling error: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respo)

}
