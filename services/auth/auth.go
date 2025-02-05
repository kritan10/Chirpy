package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github/kritan10/Chirpy/config"
	"github/kritan10/Chirpy/sql/gen"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const DURATION_ONE_HOUR = 3600 * time.Second

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func LoginHandler(apiConfig config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		body := parameters{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		user, err := apiConfig.DbQueries.GetUserByEmail(r.Context(), body.Email)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		err = CheckPasswordHash(body.Password, user.Password)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte("Incorrect email or password"))
			return
		}

		token, err := MakeJWT(user.ID, apiConfig.JwtSecret, DURATION_ONE_HOUR)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		refreshToken, err := MakeRefreshToken(apiConfig, r.Context(), user.ID)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		dat, err := json.Marshal(struct {
			ID           uuid.UUID `json:"id"`
			Email        string    `json:"email"`
			CreatedAt    time.Time `json:"created_at"`
			UpdatedAt    time.Time `json:"updated_at"`
			Token        string    `json:"token"`
			RefreshToken string    `json:"refresh_token"`
		}{
			ID:           user.ID,
			Email:        user.Email,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Token:        token,
			RefreshToken: refreshToken,
		})
		if err != nil {

		}
		w.WriteHeader(200)
		w.Write(dat)
	}
}

func RefreshTokenHandler(apiConfig config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		result, err := apiConfig.DbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		// check whether refresh token has been revoked or has expired
		if result.RevokedAt.Valid || time.Now().After(result.ExpiresAt) {
			w.WriteHeader(401)
			return
		}

		authToken, err := MakeJWT(result.UserID, apiConfig.JwtSecret, DURATION_ONE_HOUR)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		response := struct {
			Token string `json:"token"`
		}{
			Token: authToken,
		}
		dat, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		w.Write(dat)

	}
}

func RevokeRefreshTokenHandler(apiConfig config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		_, err = apiConfig.DbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
		if err != nil {
			w.WriteHeader(401)
			return
		}

		err = apiConfig.DbQueries.RevokeRefreshToken(r.Context(), refreshToken)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(204)
	}
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(expiresIn)},
		Subject:   userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Printf("error parsing jwt - %v\ntoken:%v", err, token)
		return uuid.Nil, err
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("error parsing issuer - %v", err)
		return uuid.Nil, err
	}
	return uuid.Parse(subject)
}

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	authSlice := strings.Split(authorization, " ")
	if len(authSlice) != 2 {
		return "", errors.New("invalid auth")
	}
	authType, authToken := authSlice[0], authSlice[1]
	if authType != "Bearer" {
		return "", errors.New("invalid auth")
	}

	return authToken, nil
}

func MakeRefreshToken(apiConfig config.ApiConfig, context context.Context, userId uuid.UUID) (string, error) {
	c := 32
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	token, err := apiConfig.DbQueries.CreateRefreshToken(context, gen.CreateRefreshTokenParams{
		Token:  hex.EncodeToString(b),
		UserID: userId,
	})
	if err != nil {
		return "", err
	}
	return token, nil
}
