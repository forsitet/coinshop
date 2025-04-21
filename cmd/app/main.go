package main

import (
	"coin/cmd/app/config"
	hndl "coin/internal/api/http"
	"coin/internal/database/postgres"
	"coin/internal/database/redis"
	"coin/internal/repository"
	"coin/service"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

const dbName = "coin"

func main() {
	appFlags := config.ParseFlags()
	var cfg config.AppConfig
	config.MustLoad(appFlags.ConfigPath, &cfg)
	defaultDBConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		cfg.BD.Host, cfg.BD.Port, cfg.BD.User, cfg.BD.Password)
	postgresDB, err := sql.Open("postgres", defaultDBConnStr)
	if err != nil {
		log.Fatalf("connection error to %s: %s", dbName, err)
	}

	postgres.CreateCoinRepository("coin", postgresDB)
	coinDBConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.BD.Host, cfg.BD.Port, cfg.BD.User, cfg.BD.Password, dbName)
	coinDB, err := sql.Open("postgres", coinDBConnStr)
	if err != nil {
		log.Fatalf("connection error to %s: %s", "coin", err)
	}

	postgres.InitCoinTables(coinDB)

	coinPG := postgres.NewPostgresStore(coinDB)
	log.Println(cfg.Redis.Addres)
	coinRdC := redis.NewRedisClient(cfg.Redis.Addres)
	coinRepo := repository.NewCoinRepository(coinPG, coinRdC)
	coinService := service.NewCoinService(coinRepo)
	coinHandler := hndl.NewCoinHandler(*coinService)

	r := gin.Default()
	coinHandler.WithObjectHandlers(r)

	log.Printf("Starting server on %s", cfg.HTTP.Address)
	if err := r.Run(cfg.HTTP.Address); err != nil {
		log.Fatalf("Error run server: %v", err)
	}
}
