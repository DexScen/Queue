package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	psql "github.com/DexScen/Queue/backend/internal/repository/psql"
	"github.com/DexScen/Queue/backend/internal/service"
	"github.com/DexScen/Queue/backend/internal/transport/rest"
	"github.com/DexScen/Queue/backend/pkg/database"
)

func main() {
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		Username: os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  "disable",
	})

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queuesRepo := psql.NewQueues(db)
	queuesService := service.NewQueues(queuesRepo)
	handler := rest.NewQueues(queuesService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRouter(),
	}

	log.Println("Server started at:", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
