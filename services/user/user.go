package user

import (
	"encoding/json"
	"github/kritan10/Chirpy/config"
	"github/kritan10/Chirpy/services/auth"
	"github/kritan10/Chirpy/sql/gen"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateUserHandler(apiConfig config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		decoder := json.NewDecoder(r.Body)
		body := parameters{}
		decoder.Decode(&body)

		hash, err := auth.HashPassword(body.Password)
		if err != nil {
			log.Printf("could not hash password %v", err)
			w.WriteHeader(500)
			return
		}

		user, err := apiConfig.DbQueries.CreateUser(r.Context(), gen.CreateUserParams{Email: body.Email, Password: hash})
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(500)
			return
		}
		res := struct {
			ID        uuid.UUID `json:"id"`
			Email     string    `json:"email"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		}{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		dat, err := json.Marshal(res)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(201)
		w.Header().Add("Content-Type", "application/json")
		w.Write(dat)
	}
}

func ResetUsers(apiConfig config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if apiConfig.Platform != "dev" {
			w.WriteHeader(403)
			return
		}

		err := apiConfig.DbQueries.DeleteAllUsers(r.Context())
		if err != nil {
			log.Printf("could not delete users %v", err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
	}
}
