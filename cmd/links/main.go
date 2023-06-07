package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"os"
	"project/configs"
	"project/internal/links/delivery/grpc/server"
	repositoryInMemory "project/internal/links/repository/in_memory"
	repositoryPostgres "project/internal/links/repository/postgres"
	"project/internal/links/usecase"
)

func init() {
	envPath := ".env"

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("No .env file found")
	}
}

var (
	postgres = "postgres"
	inMemory = "in-memory"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal(`В качестве хранилище укажите "postgres" или "in-memory"`)
	}

	storage := args[0]
	if storage != postgres && storage != inMemory {
		log.Fatal(`В качестве хранилище укажите "postgres" или "in-memory"`)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	yamlPath, exists := os.LookupEnv("YAML_PATH")
	if !exists {
		log.Fatal("Yaml path not found")
	}

	yamlFile, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatal(err)
	}

	var config configs.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	if storage == postgres {
		db, err := sqlx.Open(config.Postgres.DB, config.Postgres.ConnectionToDB)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err = db.Close()
			if err != nil {
				log.Error(err)
			}
		}()

		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(10)

		linksRepository := repositoryPostgres.NewLinksRepository(db)

		linksUsecase := usecase.NewLinksUsecase(linksRepository)

		grpcServer := grpc.NewServer()

		service := server.NewLinksServiceGRPCServer(grpcServer, linksUsecase)

		err = service.StartGRPCServer(config.LinksService.Addr)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		linksRepository := repositoryInMemory.NewLinksRepository()

		linksUsecase := usecase.NewLinksUsecase(linksRepository)

		grpcServer := grpc.NewServer()

		service := server.NewLinksServiceGRPCServer(grpcServer, linksUsecase)

		err = service.StartGRPCServer(config.LinksService.Addr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
