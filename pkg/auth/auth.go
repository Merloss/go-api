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
	Id       string `json:"id"`
	Username string `json:"username"`
	Role     Role   `json:"role"`
	jwt.RegisteredClaims
}

type Payload struct {
	Id       string
	Username string
	Role     Role
}

func Sign(payload *Payload, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaim{payload.Username, payload.Id, payload.Role, jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24))}})

	tokenString, err := token.SignedString(secretKey)

	return tokenString, err
}

func Decode(token string, secretKey []byte) (*CustomClaim, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &CustomClaim{}, func(t *jwt.Token) (interface{}, error) { return secretKey, nil })
	if err != nil {
		return nil, err
	}

	return jwtToken.Claims.(*CustomClaim), nil
}

func Hash(content string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(content), 10)
	errors.Must(err)

	return hash
}

func Verify(hash []byte, compareValue []byte) error {
	return bcrypt.CompareHashAndPassword(hash, compareValue)
}
