package handler

import (
	"RenderTimeEstimator/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
  	return &Handler{
    	Repository: r,
  	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/RenderServerUnits", h.GetRenderServerUnits)
	router.GET("/RenderServerUnits/:id", h.GetRenderServerUnit)
	router.GET("/AddRenderServerUnit", h.GetDraftRenderServerUnit)
	router.POST("/AddRenderServerUnit", h.AddDraftRenderServerUnit)
	router.POST("/PublishRenderServerUnit", h.PublishRenderServerUnit)
	router.POST("/DeleteRenderServerUnit", h.DeleteRenderServerUnit)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("./templates/*")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}