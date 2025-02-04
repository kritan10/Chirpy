package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func cleanBody(data string) string {
	const asterisk = "****"

	words := strings.Split(data, " ")
	cleanedBody := []string{}

	for _, word := range words {
		lowerWord := strings.ToLower(word)
		if lowerWord == "kerfuffle" || lowerWord == "sharbert" || lowerWord == "fornax" {
			cleanedBody = append(cleanedBody, asterisk)
		} else {
			cleanedBody = append(cleanedBody, word)
		}
	}
	return strings.Join(cleanedBody, " ")
}

func ValidateChirpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}

		decoder := json.NewDecoder(r.Body)
		body := parameters{}
		err := decoder.Decode(&body)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)

			res := struct {
				Error string `json:"error"`
			}{Error: "Something went wrong"}
			dat, _ := json.Marshal(res)

			w.WriteHeader(500)
			w.Write(dat)

			return
		}

		if len(body.Body) > 140 {
			log.Printf("Body length > 140")

			res := struct {
				Error string `json:"error"`
			}{Error: "Chirp is too long"}
			dat, _ := json.Marshal(res)

			w.WriteHeader(400)
			w.Write(dat)

			return
		}

		res := struct {
			CleanedBody string `json:"cleaned_body"`
		}{CleanedBody: cleanBody(body.Body)}
		dat, _ := json.Marshal(res)

		w.WriteHeader(200)
		w.Write(dat)
	}
}
