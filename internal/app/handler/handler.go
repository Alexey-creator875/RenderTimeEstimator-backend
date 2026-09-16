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

func (h *Handler) GetRenderServerUnits(ctx *gin.Context) {
	var renderServerUnits []repository.RenderServerUnit
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		renderServerUnits, err = h.Repository.GetPublishedRenderServerUnits()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		ram, err := strconv.Atoi(searchQuery)
		if err != nil {
			logrus.Error(err)
			return
		}

		renderServerUnits, err = h.Repository.GetPublishedRenderServerUnitsByRAM(ram) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "renderServerUnits.html", gin.H{
		"renderServerUnits": renderServerUnits,
		"query":  searchQuery,
	})
}

func (h *Handler) GetRenderServerUnit(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        logrus.Error(err)
        return
    }

    nextParam := ctx.Query("next")
    var renderServerUnit repository.RenderServerUnit

    switch  nextParam{
    case "":
        renderServerUnit, err = h.Repository.GetRenderServerUnit(id)
    case "true":
		renderServerUnit, err = h.Repository.GetNextPublishedRenderServerUnitTo(id)
        if err == nil {
            ctx.Redirect(http.StatusFound, "/RenderServerUnits/"+strconv.Itoa(renderServerUnit.ID))
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
        "renderServerUnit": renderServerUnit,
    })
}

func (h *Handler) GetDraftRenderServerUnit(ctx *gin.Context) {
	renderServerUnit, err := h.Repository.GetDraftRenderServerUnit()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "addRenderServerUnit.html", gin.H{
		"renderServerUnit":renderServerUnit,
	})
}