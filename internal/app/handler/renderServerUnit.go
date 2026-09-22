package handler

import (
	"RenderTimeEstimator/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetRenderServerUnits(ctx *gin.Context) {
	var renderServerUnits []ds.RenderServerUnit
	var err error

	min := ctx.Query("min")
	max := ctx.Query("max")

	if min == "" && max == "" {
		min = "8"
		max = "128"
		renderServerUnits, err = h.Repository.GetPublishedRenderServerUnits()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		min_ram, err := strconv.Atoi(min)
		if err != nil {
			logrus.Error(err)
			return
		}

		max_ram, err := strconv.Atoi(max)
		if err != nil {
			logrus.Error(err)
			return
		}

		renderServerUnits, err = h.Repository.GetPublishedRenderServerUnitsByRAM(min_ram, max_ram)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "renderServerUnits.html", gin.H{
		"renderServerUnits": renderServerUnits,
		"min":  min,
		"max": max,
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
    var renderServerUnit ds.RenderServerUnit

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

    ctx.HTML(http.StatusOK, "feedRenderServerUnits.html", gin.H{
        "renderServerUnit": renderServerUnit,
    })
}

func (h *Handler) GetDraftRenderServerUnit(ctx *gin.Context) {
	creatorID := uint(1)
 
	draftRenderServerUnit, hasDraft, err := h.Repository.GetDraftRenderServerUnit(creatorID)

	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "addRenderServerUnit.html", gin.H{
		"renderServerUnit": draftRenderServerUnit,
		"hasDraft": hasDraft,
	})
}

func (h *Handler) AddDraftRenderServerUnit(ctx *gin.Context) {
	processor := ctx.PostForm("processor")
	creatorID := uint(1)

	draftRenderServerUnit := ds.RenderServerUnit{
		Status: "draft",
		Processor: processor,
		CreatorID: creatorID,
	}

	err := h.Repository.AddDraftRenderServerUnit(draftRenderServerUnit)

	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "addRenderServerUnit.html", gin.H{
		"renderServerUnit": draftRenderServerUnit,
		"hasDraft": true,
	})
}

func (h *Handler) PublishRenderServerUnit(ctx *gin.Context) {
	coresString := ctx.PostForm("cores")
	ramString := ctx.PostForm("ram")
	description := ctx.PostForm("description")
	creatorID := uint(1)

	cores, err := strconv.Atoi(coresString)

	if err != nil {
		logrus.Error(err)
	}

	ram, err := strconv.Atoi(ramString)

	if err != nil {
		logrus.Error(err)
	}

	renderServerUnit, _, err := h.Repository.GetDraftRenderServerUnit(creatorID)

	if err != nil {
		logrus.Error(err)
	}

	renderServerUnit.Status = "published"
	renderServerUnit.Cores = cores
	renderServerUnit.RAM = ram
	renderServerUnit.Description = description

	err = h.Repository.UpdateRenderServerUnit(renderServerUnit)

	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "addRenderServerUnit.html", gin.H{
		"hasDraft": false,
	})
}
