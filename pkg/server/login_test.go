package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-api/pkg/auth"
	"go-api/pkg/entities"
	"go-api/pkg/testutil"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestLogin(t *testing.T) {
	resources := testutil.SetupResources(t, true)

	srv := New(Config{Users: resources.Users, Posts: resources.Posts})

	username := "kerim"
	password := "1234"

	buf := bytes.NewBufferString(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", buf)
	req.Header.Set("Content-Type", "application/json")

	res, err := srv.app.Test(req)
	require.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var loginRes LoginResponse

	err = json.Unmarshal(body, &loginRes)
	require.NoError(t, err)

	resources.Users.FindOne(context.TODO(), bson.D{{Key: "username", Value: username}})

	token, err := auth.Decode(loginRes.Token)
	require.NoError(t, err)

	oid, err := primitive.ObjectIDFromHex(token.Id)
	require.NoError(t, err)

	user := &entities.User{}

	err = resources.Users.FindOne(context.TODO(), bson.D{{Key: "_id", Value: oid}}).Decode(user)
	require.NoError(t, err)

	require.Equal(t, user.Username, username)
}
