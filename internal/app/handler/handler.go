package handler

import (
	"RenderTimeEstimator/internal/app/config"
	"RenderTimeEstimator/internal/app/redis"
	"RenderTimeEstimator/internal/app/repository"
	"RenderTimeEstimator/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"

	_ "RenderTimeEstimator/docs"
)

type Handler struct {
	Repository *repository.Repository
	Redis *redis.Client
	JWT config.JWTConfig
}

func NewHandler(r *repository.Repository, rd *redis.Client, jwtCfg config.JWTConfig, ) *Handler {
  	return &Handler{
    	Repository: r,
		Redis: 		rd,
		JWT:        jwtCfg,
  	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/api/sign-up", h.Register)
	router.POST("/api/sign-in", h.Login)
	router.POST("/api/sign-out", h.Logout)
	router.GET("/api/render-server-units", h.GetRenderServerUnitsAPI)
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