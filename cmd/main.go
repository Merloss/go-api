package main

import (
	"context"
	"go-api/pkg/db"
	"go-api/pkg/errors"
	"go-api/pkg/server"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := db.Connect(os.Getenv("DB_URI"), os.Getenv("DB_NAME"))
	errors.Must(err)

	srv := server.New(server.Config{Users: db.Collection("users"), Posts: db.Collection("posts")})

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		errors.Must(srv.Listen(os.Getenv("PORT")))
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	errors.Must(srv.Close(ctx))
}
