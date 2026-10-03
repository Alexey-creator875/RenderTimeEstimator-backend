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
	router.GET("/api/render-server-units", h.GetRenderServerUnitsAPI)
	router.GET("/api/render-server-units/:id", h.GetRenderServerUnitAPI)
	router.GET("/api/draft-render-server-units", h.GetDraftRenderServerUnitAPI)
	router.POST("/api/render-server-units", h.AddDraftRenderServerUnitAPI)
	router.PUT("/api/render-server-units", h.PublishRenderServerUnitAPI)
	router.DELETE("/api/render-server-units/:id", h.DeleteRenderServerUnitAPI)
	router.POST("/api/render-server-units/:id/like", h.LikeRenderServerUnitAPI)

	router.POST("/api/sign-up", h.SignUp)
	router.POST("/api/sign-in", h.SignIn)
	router.POST("/api/sign-out", h.SignOut)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("./templates/*")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"description": err.Error(),
	})
}