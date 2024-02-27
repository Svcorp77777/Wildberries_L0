package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	stan "github.com/nats-io/stan.go"

	"wildberries_L0/internal/config"
	"wildberries_L0/internal/handler"
	"wildberries_L0/internal/service"
	"wildberries_L0/internal/storage"
)

func main() {
	data, err := os.ReadFile("./configs/config.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := new(config.Config)
	if err := json.Unmarshal(data, cfg); err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	database, err := sql.Open("postgres", fmt.Sprintf("port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Base, cfg.Postgres.Sslmode))
	if err != nil {
		log.Fatal(err)
	}

	storage := storage.NewStorage(database)

	sc, err := stan.Connect("test-cluster", "sub")
	if err != nil {
		log.Fatal(err)
	}

	service := service.NewService(*storage, sc)

	err = service.RecoveryCache()
	if err != nil {
		log.Fatal(err)
	}

	err = service.Subscribe()
	if err != nil {
		log.Fatal(err)
	}

	handler := handler.NewHendler(*service, router)

	handler.Init()

	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Api.Port), router)
	if err != nil {
		log.Fatal(err)
	}

}
