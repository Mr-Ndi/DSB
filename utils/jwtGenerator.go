package utils

import (
	"dsb/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(username string, regnumber int) (string, error) {

	claims := jwt.MapClaims{
		"username":  username,
		"regnumber": regnumber,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	var jwtSecret = []byte(config.JwtKey)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(jwtSecret)

	if err != nil {
		return "", fmt.Errorf("unable to parse RSA private key")
	}
	return token.SignedString(privateKey)
}
