package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleaned_body := replaceBadWords(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned_body: cleaned_body,
	})
}

func replaceBadWords(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	listBody := strings.Split(body, " ")
	listResult := make([]string, 0)
	for _, v := range listBody {
		if containsString(v, badWords) {
			listResult = append(listResult, "****")
		} else {
			listResult = append(listResult, v)
		}
	}

	return strings.Join(listResult, " ")
}

func containsString(input string, condition []string) bool {

	for i := range condition {
		if strings.EqualFold(input, condition[i]) {
			return true
		}
	}
	return false

}
