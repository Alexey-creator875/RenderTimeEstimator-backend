package handler

import "RenderTimeEstimator/internal/app/ds"

type MessageResponse struct {
	Message string `json:"message" example:"запись успешно удалена"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"описание ошибки"`
}

type RenderServerUnitListItem struct {
	ds.RenderServerUnit
	LikesNumber int  `json:"likes_number"`
	IsMine      bool `json:"is_mine"`
}

type RenderServerUnitsResponse struct {
	RenderServerUnits []RenderServerUnitListItem `json:"renderServerUnits"`
	Min               string                     `json:"min" example:"8"`
	Max               string                     `json:"max" example:"128"`
}

type RenderServerUnitDetail struct {
	ds.RenderServerUnit
	LikesNumber int  `json:"likes_number"`
	IsLiked     bool `json:"is_liked"`
}

type RenderServerUnitDetailResponse struct {
	RenderServerUnit RenderServerUnitDetail `json:"renderServerUnit"`
}

type RenderServerUnitResponse struct {
	RenderServerUnit ds.RenderServerUnit `json:"renderServerUnit"`
}

type RenderServerUnitMessageResponse struct {
	RenderServerUnit ds.RenderServerUnit `json:"renderServerUnit"`
	Message          string              `json:"message"`
}