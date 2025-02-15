package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/Jidetireni/async-api.git/fixtures"
	"github.com/Jidetireni/async-api.git/store"
	"github.com/stretchr/testify/require"
)

func TestUserStore(t *testing.T) {

	testenv := fixtures.NewTestEnv(t)
	cleanUp := testenv.SetUpDb(t)

	// The anonymous function (func() { ... }) is used to wrap the call to cleanUp(t)
	// because t.Cleanup expects a function with no arguments and no return values
	t.Cleanup(func() {
		cleanUp(t)
	})

	now := time.Now()
	ctx := context.Background()
	userStore := store.NewUserStore(testenv.Db)
	user, err := userStore.CreateUser(ctx, "test@test123.com", "testingpasswd")
	require.NoError(t, err)
	require.Equal(t, "test@test123.com", user.Email)
	require.NoError(t, user.Validate("testingpasswd"))
	require.Less(t, now.UnixNano(), user.CreatedAt.UnixNano())

	user2, err := userStore.ByID(ctx, user.Id)
	require.NoError(t, err)
	require.Equal(t, user.Id, user2.Id)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.HashedPasswdBase64, user2.HashedPasswdBase64)
	require.Equal(t, user.CreatedAt.UnixNano(), user2.CreatedAt.UnixNano())

	user3, err := userStore.ByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.Equal(t, user.Id, user3.Id)
	require.Equal(t, user.Email, user3.Email)
	require.Equal(t, user.HashedPasswdBase64, user3.HashedPasswdBase64)
	require.Equal(t, user.CreatedAt.UnixNano(), user3.CreatedAt.UnixNano())

}
