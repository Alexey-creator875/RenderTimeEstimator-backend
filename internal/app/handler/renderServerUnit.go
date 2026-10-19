package handler

import (
	"RenderTimeEstimator/internal/app/ds"
	"fmt"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
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

// GetRenderServerUnitsAPI godoc
// @Summary      Список опубликованных серверов
// @Description  Возвращает опубликованные услуги. Фильтр по RAM: параметры min и max нужно передавать вместе, без них возвращаются все услуги. Доступно без авторизации.
// @Tags         RenderServerUnits
// @Produce      json
// @Param        min  query     int  false  "Минимальный объём RAM"
// @Param        max  query     int  false  "Максимальный объём RAM"
// @Success      200  {object}  RenderServerUnitsResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/render-server-units [get]
func (h *Handler) GetRenderServerUnitsAPI(ctx *gin.Context) {
	var renderServerUnits []ds.RenderServerUnit
	var err error

	minRamString := ctx.Query("min")
	maxRamString := ctx.Query("max")

	if minRamString == "" && maxRamString == "" {
		minRamString = "8"
		maxRamString = "128"

		renderServerUnits, err = h.Repository.GetPublishedRenderServerUnits()

		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	} else {
		minRam, err := strconv.Atoi(minRamString)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}

		maxRam, err := strconv.Atoi(maxRamString)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}

		renderServerUnits, err = h.Repository.GetPublishedRenderServerUnitsByRAM(minRam, maxRam)

		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	userId := GetUserID()

	type ExtendedRenderServerUnit struct {
		ds.RenderServerUnit
		LikesNumber int		`json:"likes_number"`
		IsMine    	bool	`json:"is_mine"`
	}

	var extendedRenderServerUnits []ExtendedRenderServerUnit

	for _, renderServerUnit := range renderServerUnits {
		likesNumber, err := h.Repository.GetLikesNumber(renderServerUnit.ID)

		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}

		isMine := renderServerUnit.CreatorID == userId

		extendedRenderServerUnit := ExtendedRenderServerUnit{
			RenderServerUnit: renderServerUnit,
			LikesNumber: likesNumber,
			IsMine: isMine,
		}

		extendedRenderServerUnits = append(extendedRenderServerUnits, extendedRenderServerUnit)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"renderServerUnits": extendedRenderServerUnits,
		"min":               minRamString,
		"max":               maxRamString,
	})
}

// GetRenderServerUnitAPI godoc
// @Summary      Получить сервер по id или следующий за ним
// @Description  Без параметра next возвращает услугу с указанным id. С next=true возвращает следующую за ней в ленте. Доступно без авторизации.
// @Tags         RenderServerUnits
// @Produce      json
// @Param        id    path      int     true   "ID услуги"
// @Param        next  query     string  false  "Если true, вернуть следующую услугу"  Enums(true)
// @Success      200   {object}  RenderServerUnitDetailResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/render-server-units/{id} [get]
func (h *Handler) GetRenderServerUnitAPI(ctx *gin.Context) {
	idString := ctx.Param("id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	nextParam := ctx.Query("next")
	var renderServerUnit ds.RenderServerUnit

	switch nextParam {
	case "":
		renderServerUnit, err = h.Repository.GetPublishedRenderServerUnit(id)

		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	case "true":
		renderServerUnit, err = h.Repository.GetPublishedRenderServerUnitNextTo(id)

		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	default:
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid parameter"))
		return
	}

	likesNumber, err := h.Repository.GetLikesNumber(renderServerUnit.ID)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	userID := GetUserID()
	isLiked, err := h.Repository.IsLikedByUser(renderServerUnit.ID, int(userID))

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	type ExtendedRenderServerUnit struct {
		ds.RenderServerUnit
		LikesNumber	int		`json:"likes_number"`
		IsLiked		bool	`json:"is_liked"`
	}

	extendedRenderServerUnit := ExtendedRenderServerUnit{
		RenderServerUnit: renderServerUnit,
		LikesNumber: likesNumber,
		IsLiked: isLiked,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"renderServerUnit": extendedRenderServerUnit,
	})
}

// GetDraftRenderServerUnitAPI godoc
// @Summary      Получить черновик текущего пользователя
// @Description  Возвращает черновик услуги, созданный текущим пользователем. Только для создателя.
// @Tags         RenderServerUnits
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  RenderServerUnitResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/draft-render-server-units [get]
func (h *Handler) GetDraftRenderServerUnitAPI(ctx *gin.Context) {
	creatorID := GetUserID()
	draftRenderServerUnit, err := h.Repository.GetDraftRenderServerUnit(creatorID)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"renderServerUnit": draftRenderServerUnit,
	})
}

// AddDraftRenderServerUnitAPI godoc
// @Summary      Создать черновик услуги
// @Description  Создаёт черновик с картинкой и видео. Создатель подставляется автоматически. Максимальный размер формы 2 МБ.
// @Tags         RenderServerUnits
// @Accept       mpfd
// @Produce      json
// @Security     ApiKeyAuth
// @Param        processor  formData  string  true   "Название процессора"
// @Param        image      formData  file    false  "Изображение (jpeg, png, gif, webp)"
// @Param        video      formData  file    false  "Видео (mp4, webm, mov, avi, mkv)"
// @Success      201  {object}  RenderServerUnitMessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/render-server-units [post]
func (h *Handler) AddDraftRenderServerUnitAPI(ctx *gin.Context) {
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
			return
		}
	} else {
		imageFound = true

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

	creatorID := GetUserID()

	draftRenderServerUnit := ds.RenderServerUnit{
		Processor: ctx.Request.FormValue("processor"),
		Status:    "draft",
		CreatorID: creatorID,
	}

	err = h.Repository.AddDraftRenderServerUnit(&draftRenderServerUnit)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if imageFound {
		err = h.Repository.AddOrReplaceRenderServerUnitImage(uint(draftRenderServerUnit.ID), imageHeader, ctx.Request.Context())

		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	if videoFound {
		err = h.Repository.AddOrReplaceRenderServerUnitVideo(uint(draftRenderServerUnit.ID), videoHeader, ctx.Request.Context());

		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"renderServerUnit": draftRenderServerUnit,
		"message":          "Черновик успешно добавлен",
	})
}

// PublishRenderServerUnitAPI godoc
// @Summary      Опубликовать черновик
// @Description  Заполняет поля черновика текущего пользователя и переводит его в статус published.
// @Tags         RenderServerUnits
// @Accept       mpfd
// @Produce      json
// @Security     ApiKeyAuth
// @Param        cores        formData  int     false  "Количество ядер"
// @Param        ram          formData  int     false  "Объём RAM"
// @Param        description  formData  string  false  "Описание"
// @Success      200  {object}  RenderServerUnitMessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/render-server-units [put]
func (h *Handler) PublishRenderServerUnitAPI(ctx *gin.Context) {
	creatorID := GetUserID()
	renderServerUnit, err := h.Repository.GetDraftRenderServerUnit(creatorID)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if renderServerUnit == nil {
		h.errorHandler(ctx, http.StatusNotFound, fmt.Errorf("черновик не найден"))
		return
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
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"renderServerUnit": renderServerUnit,
		"message":          "запись успешно обновлена",
	})
}

// DeleteRenderServerUnitAPI godoc
// @Summary      Удалить услугу
// @Description  Удаляет услугу. Доступно только её создателю.
// @Tags         RenderServerUnits
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "ID услуги"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      403  {object}  MessageResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/render-server-units/{id} [delete]
func (h *Handler) DeleteRenderServerUnitAPI(ctx *gin.Context) {
	idString := ctx.Param("id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	RenderServerUnit, err := h.Repository.GetPublishedRenderServerUnit(id)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	userID := GetUserID()

	if RenderServerUnit.CreatorID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "нет прав для удаления записи",
		})
		return
	}

	err = h.Repository.DeleteRenderServerUnit(id)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "запись успешно удалена",
	})
}


type LikeRequest struct {
	Like bool `json:"like"`
}

// LikeRenderServerUnitAPI godoc
// @Summary      Поставить или убрать лайк
// @Description  like=true ставит лайк (201), like=false снимает (200). Повторный лайк не создаёт дубликат.
// @Tags         RenderServerUnits
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id    path      int          true  "ID услуги"
// @Param        body  body      LikeRequest  true  "Поставить или снять лайк"
// @Success      200   {object}  MessageResponse  "лайк отменён"
// @Success      201   {object}  MessageResponse  "лайк поставлен"
// @Failure      400   {object}  ErrorResponse
// @Failure      403   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/render-server-units/{id}/like [post]
func (h *Handler) LikeRenderServerUnitAPI(ctx *gin.Context) {
	userID := GetUserID()

	renderServerUnitIdString := ctx.Param("id")
	renderServerUnitID, err := strconv.Atoi(renderServerUnitIdString)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var likeRequest LikeRequest
	err = ctx.ShouldBindJSON(&likeRequest)
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
				UserID:             userID,
				RenderServerUnitID: uint(renderServerUnitID),
			}

			err = h.Repository.AddLike(&like)

			if err != nil {
				h.errorHandler(ctx, http.StatusInternalServerError, err)
				return
			}
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "лайк поставлен",
		})

		return
	}

	if isLikedByUser {
		err = h.Repository.DeleteLike(userID, uint(renderServerUnitID))

		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "лайк отменён",
	})
}
