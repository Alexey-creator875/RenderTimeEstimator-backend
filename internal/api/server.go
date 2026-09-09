package api

import (
	"RenderTimeEstimator/internal/app/handler"
	"RenderTimeEstimator/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()

	r.LoadHTMLGlob("../../templates/*")
	r.Static("/static", "../../resources")

	r.GET("/RenderUnits", handler.GetRenderUnits)
	r.GET("/RenderUnits/:id", handler.GetRenderUnit)
	r.GET("/AddRenderUnit", handler.GetDraftRenderUnit)
	

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}