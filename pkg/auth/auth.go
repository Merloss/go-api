package auth

import (
	"go-api/pkg/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Role = string

const (
	ADMIN  Role = "ADMIN"
	EDITOR Role = "EDITOR"
	VIEWER Role = "VIEWER"
)

type CustomClaim struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

type Payload struct {
	Id string
}

// Sign generates a JWT (JSON Web Token) string by signing the provided payload using the specified secret key.
// It uses the HMAC SHA-256 (HS256) signing method and includes an expiration time of 24 hours in the token payload.
//
// Usage:
//
//	token, err := Sign(&Payload{Id: "user123"}, []byte("secretKey"))
func Sign(payload *Payload, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaim{payload.Id, jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24))}})

	tokenString, err := token.SignedString(secretKey)

	return tokenString, err
}

// Hash generates a bcrypt hash from the provided content using a cost factor of 10.
//
// Usage:
//
//	hashedContent := Hash("mySecretPassword")
func Hash(content string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(content), 10)
	errors.Must(err)

	return hash
}

// Verify compares a bcrypt hash with a plaintext value to check if they match.
// It uses the bcrypt.CompareHashAndPassword function for secure comparison.
//
// Usage:
//
//	err := Verify(hashedPassword, userInput)
func Verify(hash []byte, compareValue []byte) error {
	return bcrypt.CompareHashAndPassword(hash, compareValue)
}
