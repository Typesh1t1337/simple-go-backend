package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"huinya/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetJwtToken(id uint) (accessToken, RefreshToken string, err error) {
	errENV := godotenv.Load()

	if errENV != nil {
		log.Fatal("Error loading .env file")
	}

	accessClaims := jwt.MapClaims{
		"id":   id,
		"exp":  time.Now().Add(time.Hour * 48).Unix(),
		"type": "access",
	}

	accessToken, accessErr := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(os.Getenv("JWT_SECRET")))

	if accessErr != nil {
		return "", "", accessErr
	}

	refreshClaims := jwt.MapClaims{
		"id":   id,
		"exp":  time.Now().Add(time.Hour * 240).Unix(),
		"type": "refresh",
	}

	refreshToken, refreshErr := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("JWT_SECRET")))

	if refreshErr != nil {
		return "", "", refreshErr
	}

	return accessToken, refreshToken, nil
}

func GetIDFromContext(r *http.Request) (uint, bool) {
	id, ok := r.Context().Value(middleware.UserContextKey).(uint)

	return id, ok
}
