package testutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

func createMongoDB(t *testing.T) string {

	t.Helper()
	ctx := context.Background()

	mongodbContainer, err := mongodb.RunContainer(ctx)

	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, mongodbContainer.Terminate(ctx))
	})

	endpoint, err := mongodbContainer.ConnectionString(ctx)

	require.NoError(t, err)

	return endpoint
}
