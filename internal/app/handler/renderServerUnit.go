package handler

import (
	"RenderTimeEstimator/internal/app/ds"
	"fmt"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RecordLikes struct{
	LikesNumber		int
	IsLikedByUser	bool
}

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

	return slices.Contains(imageTypes, contentType)
}

func isVideo(contentType string) bool {
	videoTypes := []string{
		"video/mp4",
		"video/webm",
		"video/quicktime",
		"video/x-msvideo",
		"video/x-matroska",
	}

	return slices.Contains(videoTypes, contentType)
}

func validateFileUpload(header *multipart.FileHeader, check func(string) bool) (int, error) {
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

	if !check(contentType) {
		return http.StatusBadRequest, fmt.Errorf("некорректный файл")
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

	userId := GetUserID()
	likesMap := map[int]RecordLikes{}

	for _, renderServerUnit := range renderServerUnits {
		likesNumber, err := h.Repository.GetLikesNumber(renderServerUnit.ID)

		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}

		isLikedByUser, err := h.Repository.IsLikedByUser(renderServerUnit.ID, int(userId))

		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}

		likesMap[renderServerUnit.ID] = RecordLikes{
			LikesNumber: likesNumber,
			IsLikedByUser: isLikedByUser,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"renderServerUnits": renderServerUnits,
		"likesMap":          likesMap,
		"min":               min,
		"max":               max,
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

	switch nextParam {
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
		"likes":            likes,
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

	imageHeader, err := ctx.FormFile("image")
	imageFound := false

	if err != nil {
		if err != http.ErrMissingFile {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "File error: " + err.Error()})
			return
		}
	} else {
		imageFound = true
	}

	if imageFound {
		code, err := validateFileUpload(imageHeader, isImage)

		if err != nil {
			h.errorHandler(ctx, code, err)
			return
		}
	}

	videoHeader, err := ctx.FormFile("video")
	videoFound := false
	if err != nil {
		if err != http.ErrMissingFile {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "File error: " + err.Error()})
			return
		}
	} else {
		videoFound = true
		code, err := validateFileUpload(videoHeader, isVideo)

		if err != nil {
			h.errorHandler(ctx, code, err)
			return
		}
	}

	draftRenderServerUnit := ds.RenderServerUnit{
		Processor: ctx.Request.FormValue("processor"),
		Status: "draft",
		CreatorID: creatorID,
	}

	err = h.Repository.AddDraftRenderServerUnit(&draftRenderServerUnit)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if imageFound {
		if err = h.Repository.AddOrReplaceRenderServerUnitImage(uint(draftRenderServerUnit.ID), imageHeader, ctx); err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	if videoFound {
		if err = h.Repository.AddOrReplaceRenderServerUnitVideo(uint(draftRenderServerUnit.ID), videoHeader, ctx); err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"renderServerUnit": draftRenderServerUnit,
		"message":          "Черновик успешно добавлен",
	})
}

func (h *Handler) PublishRenderServerUnitAPI(ctx *gin.Context) {
	creatorID := GetUserID()

	renderServerUnit, err := h.Repository.GetDraftRenderServerUnit(creatorID)

	if err != nil {
		logrus.Error(err)
	}

	coresString := ctx.Request.FormValue("cores")
	if coresString != "" {
		renderServerUnit.Cores, err = strconv.Atoi(coresString)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
	}

	ramString := ctx.Request.FormValue("ram")
	if ramString != "" {
		renderServerUnit.RAM, err = strconv.Atoi(ramString)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
	}

	renderServerUnit.Status = "published"
	renderServerUnit.Description = ctx.Request.FormValue("description")

	err = h.Repository.UpdateRenderServerUnit(*renderServerUnit)

	if err != nil {
		logrus.Error(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"renderServerUnit": renderServerUnit,
		"message":          "запись успешно обновлена",
	})
}

func (h *Handler) DeleteRenderServerUnitAPI(ctx *gin.Context) {
	idString := ctx.Param("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		logrus.Error(err)
	}

	RenderServerUnit, err := h.Repository.GetPublishedRenderServerUnit(id)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
	}

	userID := GetUserID()

	if RenderServerUnit.CreatorID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message":          "нет прав для удаления записи",
		})
		return
	}

	err = h.Repository.DeleteRenderServerUnit(id)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":          "запись успешно удалена",
	})
}

type LikeRequest struct {
	Like bool `json:"like"`
}

func (h *Handler) LikeRenderServerUnitAPI(ctx *gin.Context) {
	userID := GetUserID()

	renderServerUnitIdString := ctx.Param("id")
	renderServerUnitID, err := strconv.Atoi(renderServerUnitIdString)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var likeRequest LikeRequest
	err = ctx.ShouldBindJSON(&likeRequest);
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	isLikedByUser, err := h.Repository.IsLikedByUser(renderServerUnitID, int(userID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if likeRequest.Like {
		if !isLikedByUser {
			like := ds.Likes{
				UserID: userID,
				RenderServerUnitID: uint(renderServerUnitID),
			}
			h.Repository.AddLike(&like)
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message":          "лайк поставлен",
		})

		return
	}

	if isLikedByUser {
		h.Repository.DeleteLike(userID, uint(renderServerUnitID))
	}

	ctx.JSON(http.StatusOK, gin.H{
			"message":          "лайк отменён",
		})
}
