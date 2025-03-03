package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sant470/trademark/config"
)

func GenerateJWT(username string, role string) (string, error) {
	// Define token claims (payload)
	secret := config.GetAppConfig("config.yaml", ".").JWT
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":      time.Now().Unix(),                     // Issued At
	}
	// Create the token with claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and generate the token string
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
