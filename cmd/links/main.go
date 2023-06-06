package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	goRedis "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"os"
	"project/configs"
	"project/internal/links/delivery/grpc/server"
	repositoryPostgres "project/internal/links/repository/postgres"
	repositoryRedis "project/internal/links/repository/redis"
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
	redis    = "redis"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal(`В качестве хранилище укажите "postgres" или "redis"`)
	}

	storage := args[0]
	if storage != postgres && storage != redis {
		log.Fatal(`В качестве хранилище укажите "postgres" или "redis"`)
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
		db := goRedis.NewClient(&goRedis.Options{
			Addr: config.Redis.Addr,
		})
		defer func() {
			err = db.Close()
			if err != nil {
				log.Error(err)
			}
		}()

		linksRepository := repositoryRedis.NewLinksRepository(db)

		linksUsecase := usecase.NewLinksUsecase(linksRepository)

		grpcServer := grpc.NewServer()

		service := server.NewLinksServiceGRPCServer(grpcServer, linksUsecase)

		err = service.StartGRPCServer(config.LinksService.Addr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
