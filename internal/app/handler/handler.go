package handler

import (
	"RenderTimeEstimator/internal/app/config"
	"RenderTimeEstimator/internal/app/repository"
	"RenderTimeEstimator/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
	JWT config.JWTConfig
}

func NewHandler(r *repository.Repository, jwtCfg config.JWTConfig) *Handler {
  	return &Handler{
    	Repository: r,
		JWT:        jwtCfg,
  	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.POST("/api/sign-up", h.Register)
	router.POST("/api/sign-in", h.Login)
	router.POST("/api/sign-out", h.SignOut)

	router.Use(h.WithAuthCheck(role.Guest, role.Admin)).GET("/api/render-server-units", h.GetRenderServerUnitsAPI)
	router.GET("/api/render-server-units/:id", h.GetRenderServerUnitAPI)

	router.Use(h.WithAuthCheck(role.Admin)).GET("/api/draft-render-server-units", h.GetDraftRenderServerUnitAPI)
	router.POST("/api/render-server-units", h.AddDraftRenderServerUnitAPI)
	router.PUT("/api/render-server-units", h.PublishRenderServerUnitAPI)
	router.DELETE("/api/render-server-units/:id", h.DeleteRenderServerUnitAPI)
	router.POST("/api/render-server-units/:id/like", h.LikeRenderServerUnitAPI)

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