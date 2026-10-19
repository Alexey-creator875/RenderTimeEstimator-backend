package main

import (
	"RenderTimeEstimator/internal/app/config"
	"RenderTimeEstimator/internal/app/dsn"
	"RenderTimeEstimator/internal/app/handler"
	"RenderTimeEstimator/internal/app/redis"
	"RenderTimeEstimator/internal/app/repository"
	"RenderTimeEstimator/internal/pkg"
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title BITOP
// @version 1.0
// @description Bmstu Open IT Platform

// @contact.name API Support
// @contact.url https://vk.com/bmstu_schedule
// @contact.email bitop@spatecon.ru

// @license.name AS IS (NO WARRANTY)

// @host 127.0.0.1
// @schemes https http
// @BasePath /

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.New(&repository.RepositorySettings{
		PostgresDSN: postgresString,
		MinioEndpoint: os.Getenv("MINIO_ENDPOINT"),
		MinioAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey: os.Getenv("MINIO_SECRET_KEY"),
		MinioBucketName: os.Getenv("MINIO_BUCKET_NAME"),	
	})

	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	ctx := context.Background()

	rd, err := redis.New(ctx, conf.Redis)
	if err != nil {
		logrus.Fatalf("error initializing redis: %v", err)
	}
	defer rd.Close()

	hand := handler.NewHandler(rep, rd, conf.JWT)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}