package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-api/pkg/entities"
	"go-api/pkg/testutil"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	resources := testutil.SetupResources(t, true)

	srv := New(Config{Users: resources.Users, Posts: resources.Posts})

	title := "testtitle"
	description := "test post description"

	buf := bytes.NewBufferString(fmt.Sprintf(`{"title": "%s", "description": "%s"}`, title, description))
	req := httptest.NewRequest(http.MethodPost, "/api/posts", buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resources.AdminUser.Token))

	res, err := srv.app.Test(req)
	require.NoError(t, err)

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var createPostResponse CreatePostResponse

	err = json.Unmarshal(body, &createPostResponse)
	require.NoError(t, err)

	require.Equal(t, title, createPostResponse.Post.Title)
	require.Equal(t, description, createPostResponse.Post.Description)
	require.Equal(t, entities.PENDING, createPostResponse.Post.Status)
}
