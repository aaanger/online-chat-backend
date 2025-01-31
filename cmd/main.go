package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"onlineChat/internal/routes"
	"onlineChat/internal/users"
	"onlineChat/internal/ws"
	"onlineChat/pkg/db"
	"onlineChat/pkg/redis"
	"os"
)

type Server struct {
	httpServer *http.Server
}

func (srv *Server) Run(port string, handler http.Handler) error {
	srv.httpServer = &http.Server{
		Addr:    port,
		Handler: handler,
	}
	return srv.httpServer.ListenAndServe()
}

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

	redisCfg := redis.RedisConfig{
		Address:  os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}

	hub := ws.NewHub(redisCfg)
	go hub.Run()

	userRepo := users.NewUserRepository(db)
	userService := users.NewUserService(userRepo)
	userHandler := users.NewUserHandler(userService)

	chatRepo := ws.NewChatRepository(db)
	chatService := ws.NewChatService(chatRepo)
	chatHandler := ws.NewChatHandler(hub, chatService)

	srv := new(Server)
	err = srv.Run(":3000", routes.PathHandler(userHandler, chatHandler))
	if err != nil {
		log.Fatalf("error running the server: %w", err)
	}
}
