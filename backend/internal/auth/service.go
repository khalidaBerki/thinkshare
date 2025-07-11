package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(getSecret())

func getSecret() string {
	secret := os.Getenv("JWT_SECRET")
	return secret
}

func GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,                                // Payload : ID de l'utilisateur
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Expiration du token : 24h
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseJWT(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token") // token expiré ou signature invalide
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("cannot parse claims")
	}

	userID, ok := claims["user_id"].(float64) // json number = float64
	if !ok {
		return 0, errors.New("user_id not found in token")
	}

	return int(userID), nil
}
