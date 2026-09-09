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

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		renderUnits, err = h.Repository.GetPublishedRenderUnits()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		renderUnits, err = h.Repository.GetRenderUnitsByProcessor(searchQuery) // в ином случае ищем заказ по заголовку
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
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        logrus.Error(err)
        return
    }

    nextParam := ctx.Query("next")
    var renderUnit repository.RenderUnit

    switch  nextParam{
    case "":
        renderUnit, err = h.Repository.GetRenderUnit(id)
    case "true":
		renderUnit, err = h.Repository.GetNextPublishedRenderUnitTo(id)
        if err == nil {
            ctx.Redirect(http.StatusFound, "/RenderUnits/"+strconv.Itoa(renderUnit.ID))
            return
        }
    default:
        logrus.Error(err)
        ctx.String(http.StatusBadRequest, "параметр next должен быть true или отсутствовать")
        return
    }

    if err != nil {
        logrus.Error(err)
        ctx.String(http.StatusNotFound, "запись не найдена")
        return
    }

    ctx.HTML(http.StatusOK, "video.html", gin.H{
        "renderUnit": renderUnit,
    })
}

func (h *Handler) GetDraftRenderUnit(ctx *gin.Context) {
	renderUnit, err := h.Repository.GetDraftRenderUnit()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "addRenderUnit.html", gin.H{
		"renderUnit":renderUnit,
	})
}

// func (h *Handler) GetNextRenderUnit(ctx *gin.Context) {
// 	idStr := ctx.Param("id")
// 	id, err := strconv.Atoi(idStr)

// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	next := ctx.Param("next")

// 	if next != "true" {
// 		ctx.HTML(http.StatusNotFound, "", gin.H{})
// 	}

// 	renderUnit, err := h.Repository.GetNextPublishedRenderUnitTo(id)
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	ctx.HTML(http.StatusOK, "video.html", gin.H{
// 		"renderUnit":renderUnit,
// 	})
// }