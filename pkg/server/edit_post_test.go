package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-api/pkg/testutil"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEditPost(t *testing.T) {
	resources := testutil.SetupResources(t, true)

	srv := New(Config{Users: resources.Users, Posts: resources.Posts})

	title := "edittesttitle"
	description := "edit test post description"

	buf := bytes.NewBufferString(fmt.Sprintf(`{"title": "%s", "description": "%s"}`, title, description))
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/posts/%s", resources.PostApproved.Id), buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resources.EditorUser.Token))

	res, err := srv.app.Test(req)
	require.NoError(t, err)

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var createPostResponse CreatePostResponse

	err = json.Unmarshal(body, &createPostResponse)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
}
