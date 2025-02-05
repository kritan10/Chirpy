package config

import "github/kritan10/Chirpy/sql/gen"

type ApiConfig struct {
	Platform  string
	DbQueries *gen.Queries
	JwtSecret string
}
