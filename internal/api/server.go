package api

import (
  "github.com/gin-gonic/gin"
  "github.com/sirupsen/logrus"
  "log"
  "RenderTimeEstimator/internal/app/handler"
  "RenderTimeEstimator/internal/app/repository"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	// добавляем наш html/шаблон
	r.LoadHTMLGlob("../../templates/*")

	r.GET("/hello", handler.GetOrders)
	r.GET("/order/:id", handler.GetOrder) // вот наш новый обработчик

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}
