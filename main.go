package main

import (
	"database/sql"
	"github/kritan10/Chirpy/config"
	"github/kritan10/Chirpy/services/auth"
	"github/kritan10/Chirpy/services/chirp"
	"github/kritan10/Chirpy/services/user"
	"github/kritan10/Chirpy/sql/gen"
	"github/kritan10/Chirpy/utils"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Printf("db error %v", err)
	}

	dbQueries := gen.New(db)

	apiConfig := config.ApiConfig{
		Platform:  os.Getenv("PLATFORM"),
		DbQueries: dbQueries,
		JwtSecret: os.Getenv("JWT_SECRET"),
	}
	metrics := utils.MetricsConfig{}

	mux := http.ServeMux{}

	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", metrics.MiddlewareMetricsInc(fileServer))

	mux.HandleFunc("GET /api/healthz", utils.HealthzHandler())

	mux.HandleFunc("POST /api/login", auth.LoginHandler(apiConfig))
	mux.HandleFunc("POST /api/refresh", auth.RefreshTokenHandler(apiConfig))
	mux.HandleFunc("POST /api/revoke", auth.RevokeRefreshTokenHandler(apiConfig))

	mux.HandleFunc("POST /api/users", user.CreateUserHandler(apiConfig))
	mux.HandleFunc("POST /api/users/reset", user.ResetUsers(apiConfig))

	mux.HandleFunc("POST /api/chirps", chirp.CreateChirpHandler(apiConfig))
	mux.HandleFunc("GET /api/chirps", chirp.GetAllChirpsHandler(apiConfig))
	mux.HandleFunc("GET /api/chirps/{id}", chirp.GetChirpByIdHandler(apiConfig))

	mux.HandleFunc("GET /admin/metrics", metrics.MetricsHandler())
	mux.HandleFunc("POST /admin/reset", user.ResetUsers(apiConfig))

	server := http.Server{Addr: ":8080", Handler: &mux}
	server.ListenAndServe()
}
