package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
	"os"
	"project/configs"
	"project/internal/links/delivery/grpc/client"
	"project/internal/links/delivery/http"
	myMiddleware "project/internal/middleware"
)

func init() {
	envPath := ".env"

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetReportCaller(true)

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

	grpcConnLinks, err := grpc.Dial(
		config.LinksService.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("cant connect to grpc ", err)
	}
	defer func() {
		err = grpcConnLinks.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	linksService := client.NewLinksServiceGRPSClient(grpcConnLinks)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods:     config.Cors.AllowMethods,
		AllowOrigins:     config.Cors.AllowOrigins,
		AllowCredentials: config.Cors.AllowCredentials,
		AllowHeaders:     config.Cors.AllowHeaders,
	}))

	e.Use(myMiddleware.LoggerMiddleware)

	http.NewLinksHandler(e, linksService)

	e.Logger.Fatal(e.Start(config.Server.Port))
}
