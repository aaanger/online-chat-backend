package cmd

import (
	"github.com/joho/godotenv"
	"log"
	"onlineChat/db"
	"onlineChat/handler"
	"onlineChat/internal/users"
	"onlineChat/internal/ws"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env: %w", err)
	}

	db, err := db.Open(db.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		Username: os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DBName:   os.Getenv("PSQL_DBNAME"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	})
	if err != nil {
		log.Fatalf("Error opening db: %w", err)
	}

	defer db.Close()

	userRepo := users.NewRepository(db)
	userService := users.NewUserService(userRepo)
	userHandler := users.NewHandler(userService)

	hub := ws.NewHub()
	go hub.Run()

	srv := new(Server)
	err = srv.Run(":3000", handler.PathHandler(userHandler))
	if err != nil {
		log.Fatalf("error running the server: %w", err)
	}
}
