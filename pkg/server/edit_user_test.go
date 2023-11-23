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
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEditUser(t *testing.T) {
	resources := testutil.SetupResources(t, true)

	srv := New(Config{Users: resources.Users, Posts: resources.Posts})

	username := "newUsername"
	roles := strings.Join([]string{fmt.Sprintf(`"%s"`, string(entities.VIEWER)), fmt.Sprintf(`"%s"`, string(entities.EDITOR))}, ", ")

	buf := bytes.NewBufferString(fmt.Sprintf(`{"username": "%s", "roles": [%s]}`, username, roles))
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/users/%s/edit", resources.ViewerUser.User.Id), buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resources.AdminUser.Token))

	res, err := srv.app.Test(req)
	require.NoError(t, err)

	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	var updateUserResponse UpdateUserResponse

	err = json.Unmarshal(body, &updateUserResponse)
	require.NoError(t, err)

}
