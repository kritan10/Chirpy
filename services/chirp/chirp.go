package chirp

import (
	"encoding/json"
	"github/kritan10/Chirpy/config"
	"github/kritan10/Chirpy/sql/gen"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func CreateChirpHandler(apiConfig config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type parameters struct {
			Body   string    `json:"body"`
			UserId uuid.UUID `json:"user_id"`
		}

		decoder := json.NewDecoder(r.Body)
		body := parameters{}
		decoder.Decode(&body)

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

		chirp, err := apiConfig.DbQueries.CreateChirp(r.Context(), gen.CreateChirpParams{Body: cleanChirpBody(body.Body), UserID: body.UserId})
		if err != nil {
			log.Printf("Error creating chirp: %v", err)
			w.WriteHeader(500)
			return
		}

		res := struct {
			ID        uuid.UUID
			CreatedAt time.Time
			UpdatedAt time.Time
			Body      string
			UserID    uuid.UUID
		}{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}

		dat, err := json.Marshal(res)
		if err != nil {
			log.Printf("Error encoding chirp response: %v", err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(201)
		w.Write(dat)
	}
}

func GetAllChirpsHandler(apiConfig config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := apiConfig.DbQueries.GetChirps(r.Context())
		if err != nil {

		}

		type ChirpsJSON struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Body      string    `json:"body"`
			UserID    uuid.UUID `json:"user_id"`
		}

		chirpsJson := []ChirpsJSON{}
		for _, chirp := range chirps {
			chirpsJson = append(chirpsJson, ChirpsJSON{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			})
		}
		dat, err := json.Marshal(chirpsJson)

		if err != nil {

		}

		w.WriteHeader(200)
		w.Write(dat)
	}
}

func GetChirpByIdHandler(apiConfig config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		chirpId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			log.Printf("could not parse uuid %v", err)
			w.WriteHeader(400)
			// w.Write()
		}

		chirp, err := apiConfig.DbQueries.GetChirpById(r.Context(), chirpId)
		if err != nil {

		}

		dat, err := json.Marshal(struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Body      string    `json:"body"`
			UserID    uuid.UUID `json:"user_id"`
		}{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})

		if err != nil {

		}

		w.WriteHeader(200)
		w.Write(dat)
	}
}

func cleanChirpBody(data string) string {
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
