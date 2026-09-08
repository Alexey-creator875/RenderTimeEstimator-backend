package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type RenderUnit struct {
  	ID int
  	Title string
	GPU string
	RAM string
	Storage string
	ImagePath string
}

func (r *Repository) GetRenderUnits() ([]RenderUnit, error) {
	renderUnits := []RenderUnit{
		{
			ID:    1,
			Title: "2080 Ti",
			GPU: "2080 Ti",
			RAM: "32768 МБ ОЗУ",
			Storage: "160 ГБ SSD",
			ImagePath: "2080-Ti.jpg",
		},
		{
			ID:    2,
			Title: "3090",
			GPU: "3090",
			RAM: "32768 МБ ОЗУ",
			Storage: "160 ГБ SSD",
			ImagePath: "2080-Ti.jpg",
		},
		{
			ID:    3,
			Title: "4090",
			GPU: "4090",
			RAM: "32768 МБ ОЗУ",
			Storage: "160 ГБ SSD",
			ImagePath: "2080-Ti.jpg",
		},
	}

	if len(renderUnits) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return renderUnits, nil
}

func (r *Repository) GetRenderUnit(id int) (RenderUnit, error) {
	renderUnits, err := r.GetRenderUnits()
	if err != nil {
		return RenderUnit{}, err
	}

	for _, renderUnit := range renderUnits {
		if renderUnit.ID == id {
			return renderUnit, nil
		}
	}
	return RenderUnit{}, fmt.Errorf("заказ не найден")
}

func (r *Repository) GetRenderUnitsByTitle(title string) ([]RenderUnit, error) {
	renderUnits, err := r.GetRenderUnits()
	if err != nil {
		return []RenderUnit{}, err
	}

	var result []RenderUnit
	for _, renderUnit := range renderUnits {
		if strings.Contains(strings.ToLower(renderUnit.Title), strings.ToLower(title)) {
			result = append(result, renderUnit)
		}
	}

	return result, nil
}