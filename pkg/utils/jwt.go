package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWTToken(userId uuid.UUID, secretKey string) (string, error) {
	jwtExp := os.Getenv("JWT_EXPIRED")

	// convert jwtexp string to int
	jwtExpInt, err := strconv.Atoi(jwtExp)
	if err != nil {
		log.Println("Error converting JWT expiration time to integer:", err)
		return "", err
	}

	// create a new token with claims and secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId.String(),
		"exp":     time.Now().Add(time.Hour * time.Duration(jwtExpInt)).Unix(),
	})

	// sign the token with secret key and return it
	return token.SignedString([]byte(secretKey))
}

func ValidateJWTToken(tokenString string, secretKey string) (*jwt.MapClaims, error) {
	// parse the token using the secret key
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})

	// handle any errors that occur during parsing
	if err != nil {
		log.Printf("Error validating JWT token: %v", err)
		return nil, err
	}

	// check if the token is valid and has not expired
	if !token.Valid {
		log.Printf("JWT token is not valid")
		return nil, fmt.Errorf("JWT token is not valid")
	}

	// check if the token is valid and has not expired
	// Safe to assert as *jwt.MapClaims
	if claims, ok := token.Claims.(*jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid claims type")
}
