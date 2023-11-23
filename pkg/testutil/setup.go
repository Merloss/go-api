package testutil

import (
	"context"
	"go-api/pkg/auth"
	"go-api/pkg/db"
	"go-api/pkg/entities"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type resources struct {
	Users        *mongo.Collection
	Posts        *mongo.Collection
	AdminUser    User
	EditorUser   User
	ViewerUser   User
	PostPending  *entities.Post
	PostApproved *entities.Post
}

type User struct {
	User  *entities.User
	Token string
}

func SetupResources(t *testing.T, mock bool) resources {
	t.Helper()
	t.Setenv("JWT_SECRET", "aLuK26TNAzvIRdXwuFTZPBKW9At8rL6s")

	uri, dbName := os.Getenv("DB_URI"), os.Getenv("DB_NAME")
	if uri == "" && dbName == "" {
		uri = createMongoDB(t)
		dbName = "test"
	}

	database, err := db.Connect(uri, dbName)
	require.NoError(t, err)

	r := resources{Users: database.Collection("users"), Posts: database.Collection("posts")}

	if mock {
		r.ViewerUser = CreateUser(t, r, &entities.User{Username: "kerim", Password: string(auth.Hash("1234")), Roles: []auth.Role{auth.VIEWER}})
		r.EditorUser = CreateUser(t, r, &entities.User{Username: "kerimeditor", Password: string(auth.Hash("1234")), Roles: []auth.Role{auth.EDITOR, auth.VIEWER}})
		r.AdminUser = CreateUser(t, r, &entities.User{Username: "kerimadmin", Password: string(auth.Hash("1234")), Roles: []auth.Role{auth.ADMIN, auth.VIEWER}})

		title := "testposttitle"
		description := "test post description"

		r.PostApproved = CreatePost(t, r, &entities.Post{Title: title, Description: description, Status: entities.APPROVED})
		r.PostPending = CreatePost(t, r, &entities.Post{Title: title, Description: description, Status: entities.PENDING})
	}

	return r
}

func CreateUser(t *testing.T, r resources, user *entities.User) User {

	res, err := r.Users.InsertOne(context.TODO(), user)
	require.NoError(t, err)

	user.Id = res.InsertedID.(primitive.ObjectID).Hex()

	token, err := auth.Sign(&auth.Payload{Id: user.Id}, []byte(os.Getenv("JWT_SECRET")), nil)
	require.NoError(t, err)

	return User{User: user, Token: token}
}

func CreatePost(t *testing.T, r resources, post *entities.Post) *entities.Post {
	res, err := r.Posts.InsertOne(context.TODO(), post)
	require.NoError(t, err)

	post.Id = res.InsertedID.(primitive.ObjectID).Hex()

	return post
}
