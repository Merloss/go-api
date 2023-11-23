package server

import (
	"fmt"
	"go-api/pkg/testutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeletePost(t *testing.T) {
	resources := testutil.SetupResources(t, true)

	srv := New(Config{Users: resources.Users, Posts: resources.Posts})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/posts/%s", resources.PostApproved.Id), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resources.AdminUser.Token))

	res, err := srv.app.Test(req)
	require.NoError(t, err)

	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
}
