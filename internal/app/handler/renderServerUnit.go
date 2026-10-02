package handler

import (
	"RenderTimeEstimator/internal/app/ds"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetUserID() uint {
	return uint(1)
}

func isImage(contentType string) bool {
	imageTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, t := range imageTypes {
		if contentType == t {
			return true
		}
	}
	return false
}

func validateFileUpload(header *multipart.FileHeader) (int, error) {
	// Окрываем чтение файлового потока
	file, err := header.Open()

	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("не удалось получить файл")
	}

	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)

	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("не удалось прочитать файл")
	}

	contentType := http.DetectContentType(buffer)

	if !isImage(contentType) {
		return http.StatusBadRequest, fmt.Errorf("файл должен быть изображением")
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("ошибка при обработке файла")
	}

	return 0, nil
}


func (h *Handler) GetRenderServerUnitsAPI(ctx *gin.Context) {
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

	likesMap := map[int]int{}

	for _, renderServerUnit := range renderServerUnits {
		likesNumber, err := h.Repository.GetLikesNumber(renderServerUnit.ID)

		if err != nil {
			return
		}

		likesMap[renderServerUnit.ID] = likesNumber
	}

	ctx.JSON(http.StatusOK, gin.H{
		"renderServerUnits": renderServerUnits,
		"likesMap": likesMap,
		"min":  min,
		"max": max,
	})
}

func (h *Handler) GetRenderServerUnitAPI(ctx *gin.Context) {
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
        renderServerUnit, err = h.Repository.GetPublishedRenderServerUnit(id)
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

	likes, err := h.Repository.GetLikesNumber(renderServerUnit.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"renderServerUnit": renderServerUnit,
		"likes": likes,
	})
}

func (h *Handler) GetDraftRenderServerUnitAPI(ctx *gin.Context) {
	creatorID := GetUserID()
 
	draftRenderServerUnit, err := h.Repository.GetDraftRenderServerUnit(creatorID)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"renderServerUnit": draftRenderServerUnit,
	})
}

func (h *Handler) AddDraftRenderServerUnitAPI(ctx *gin.Context) {
	creatorID := GetUserID()

	err := ctx.Request.ParseMultipartForm(2 << 20)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	header, err := ctx.FormFile("image")

	fileFound := false

	if err != nil {
		if err != http.ErrMissingFile {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "File error: " + err.Error()})
			return
		}
	} else {
		fileFound = true
	}

	if fileFound {
		code, err := validateFileUpload(header)

		if err != nil {
			h.errorHandler(ctx, code, err)
			return
		}
	}

	draftRenderServerUnit := ds.RenderServerUnit{
		Processor:	ctx.Request.FormValue("processor"),
		CreatorID:	creatorID,
	}

	err = h.Repository.AddDraftRenderServerUnit(&draftRenderServerUnit)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if fileFound {
		if err = h.Repository.AddOrReplaceRenderServerUnitImage(uint(draftRenderServerUnit.ID), header, ctx); err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"renderServerUnit": draftRenderServerUnit,
		"message": "Черновик успешно добавлен",
	})
}

func (h *Handler) PublishRenderServerUnitAPI(ctx *gin.Context) {
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

	renderServerUnit, err := h.Repository.GetDraftRenderServerUnit(creatorID)

	if err != nil {
		logrus.Error(err)
	}

	renderServerUnit.Status = "published"
	renderServerUnit.Cores = cores
	renderServerUnit.RAM = ram
	renderServerUnit.Description = description

	err = h.Repository.UpdateRenderServerUnit(*renderServerUnit)

	if err != nil {
		logrus.Error(err)
	}

	ctx.JSON(http.StatusOK, gin.H{})

	// ctx.HTML(http.StatusOK, "addRenderServerUnit.html", gin.H{
	// 	"hasDraft": false,
	// })
}

func (h *Handler) DeleteRenderServerUnitAPI(ctx *gin.Context) {
	idString := ctx.PostForm("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		logrus.Error(err)
	}

	err = h.Repository.DeleteRenderServerUnit(id)

	if err != nil {
		logrus.Error(err)
	}

	ctx.JSON(http.StatusOK, gin.H{})

	// ctx.Redirect(http.StatusFound, "/RenderServerUnits")
}
