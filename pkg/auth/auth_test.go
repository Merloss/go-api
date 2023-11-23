package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignAndDecodeToken(t *testing.T) {
	secretKey := []byte("secret")
	payload := &Payload{Id: "123456789"}

	tokenString, err := Sign(payload, secretKey, nil)
	require.NoError(t, err)

	decodedClaim, err := Decode(tokenString)
	require.NoError(t, err)

	require.Equal(t, payload.Id, decodedClaim.Id)
}

func TestHashAndVerifyPassword(t *testing.T) {
	password := "t0pS3cReTPasswOrd"

	hashedPassword := Hash(password)
	require.NotEmpty(t, hashedPassword)

	err := Verify(hashedPassword, []byte(password))
	require.NoError(t, err)
}

func TestVerifyInvalidPassword(t *testing.T) {
	correctPassword := "t0pS3cReTPasswOrd"
	incorrectPassword := "wrongPassword"

	hashedPassword := Hash(correctPassword)

	err := Verify(hashedPassword, []byte(incorrectPassword))
	require.Error(t, err)
}
