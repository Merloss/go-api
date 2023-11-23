package server

import (
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

func TestGetPendingPosts(t *testing.T) {
	resources := testutil.SetupResources(t, true)

	srv := New(Config{Users: resources.Users, Posts: resources.Posts})

	req := httptest.NewRequest(http.MethodGet, "/api/posts/pending", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resources.AdminUser.Token))

	res, err := srv.app.Test(req)
	require.NoError(t, err)

	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var postsResponse PostsResponse

	err = json.Unmarshal(body, &postsResponse)
	require.NoError(t, err)

	for _, pr := range postsResponse.Posts {
		require.Equal(t, entities.PENDING, pr.Status)
	}

}
