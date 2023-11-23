package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-api/pkg/errors"
	"go-api/pkg/testutil"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	resources := testutil.SetupResources(t, false)

	srv := New(Config{Users: resources.Users, Posts: resources.Posts})

	username := "kerim"
	password := "1234"

	buf := bytes.NewBufferString(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", buf)
	req.Header.Set("Content-Type", "application/json")

	res, err := srv.app.Test(req)
	require.NoError(t, err)

	defer res.Body.Close()

	require.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var registerResponse RegisterResponse

	err = json.Unmarshal(body, &registerResponse)
	require.NoError(t, err)

	require.Equal(t, username, registerResponse.User.Username)
}

func TestRegisterIsUsernameAlreadyTaken(t *testing.T) {
	resources := testutil.SetupResources(t, true)

	srv := New(Config{Users: resources.Users, Posts: resources.Posts})

	username := "kerim"
	password := "1234"

	buf := bytes.NewBufferString(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", buf)
	req.Header.Set("Content-Type", "application/json")

	res, err := srv.app.Test(req)
	require.NoError(t, err)

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var errResponse errors.ErrorResponse

	err = json.Unmarshal(body, &errResponse)
	require.NoError(t, err)

	require.Equal(t, res.StatusCode, http.StatusBadRequest)
	require.Equal(t, errResponse.Error.Message, "username is already taken")
	require.Equal(t, errResponse.Error.Code, errors.BAD_REQUEST)
}
