package handler

import (
	"RenderTimeEstimator/internal/app/repository"
	"net/http"
	"strconv"

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

func (h *Handler) GetRenderUnits(ctx *gin.Context) {
	var renderUnits []repository.RenderUnit
	var err error

	searchQuery := ctx.Query("query") // получаем значение из поля поиска
	if searchQuery == "" {            // если поле поиска пусто, то просто получаем из репозитория все записи
		renderUnits, err = h.Repository.GetRenderUnits()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		renderUnits, err = h.Repository.GetRenderUnitsByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "renderUnits.html", gin.H{
		"renderUnits": renderUnits,
		"query":  searchQuery, // передаем введенный запрос обратно на страницу
		// в ином случае оно будет очищаться при нажатии на кнопку
	})
}

func (h *Handler) GetRenderUnit(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	renderUnit, err := h.Repository.GetRenderUnit(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "video.html", gin.H{
		"renderUnit":renderUnit,
	})
}