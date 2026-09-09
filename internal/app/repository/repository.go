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
	Status string
	Processor string
	Cores int
	RAM int
	Description string
	Likes int
	ImagePath string
	VideoPath string
}

func (r *Repository) GetRenderUnits() ([]RenderUnit, error) {
	renderUnits := []RenderUnit{
		{
			ID:    1,
			Status: "published",
			Processor: "2080 Ti",
			Cores: 20,
			RAM: 16,
			Description: "Надёжный выбор для рендеринга и игр. Отличная производительность в 2K и 4K, поддерживает все современные технологии трассировки лучей. Идеально для студийных проектов средней сложности.",
			Likes: 125,
			ImagePath: "2080-Ti.jpg",
			VideoPath: "2080-Ti.mp4",
		},
		{
			ID:    2,
			Status: "published",
			Processor: "3090",
			Cores: 20,
			RAM: 24,
			Description: "Флагманская карта с 24 ГБ видеопамяти. Мгновенный рендеринг сложных сцен, работа с нейросетями и 8K-видео. Профессиональный инструмент для 3D-художников и дата-сайентистов.",
			Likes: 243,
			ImagePath: "3090.jpg",
			VideoPath: "2080-Ti.mp4",
		},
		{
			ID:    3,
			Status: "published",
			Processor: "4090",
			Cores: 20,
			RAM: 24,
			Description: "Абсолютный лидер производительности. Справляется с самыми тяжёлыми задачами: симуляция физики, объёмная визуализация, рендеринг в реальном времени. Будущее уже сегодня.",
			Likes: 377,
			ImagePath: "4090.png",
			VideoPath: "2080-Ti.mp4",
		},
		{
			ID:    4,
			Status: "draft",
			Processor: "",
			Cores: 0,
			RAM: 0,
			Description: "Описание",
			Likes: 0,
			ImagePath: "4090.png",
			VideoPath: "2080-Ti.mp4",
		},
		{
			ID:    5,
			Status: "deleted",
			Processor: "4090",
			Cores: 20,
			RAM: 24,
			Description: "Абсолютный лидер производительности. Справляется с самыми тяжёлыми задачами: симуляция физики, объёмная визуализация, рендеринг в реальном времени. Будущее уже сегодня.",
			Likes: 17,
			ImagePath: "4090.png",
			VideoPath: "2080-Ti.mp4",
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

func (r *Repository) GetRenderUnitsByProcessor(processor string) ([]RenderUnit, error) {
	renderUnits, err := r.GetRenderUnits()
	if err != nil {
		return []RenderUnit{}, err
	}

	var result []RenderUnit
	for _, renderUnit := range renderUnits {
		if strings.Contains(strings.ToLower(renderUnit.Processor), strings.ToLower(processor)) {
			result = append(result, renderUnit)
		}
	}

	return result, nil
}

func (r *Repository) GetPublishedRenderUnits() ([]RenderUnit, error) {
	renderUnits, err := r.GetRenderUnits()
	if err != nil {
		return []RenderUnit{}, err
	}

	var result []RenderUnit
	for _, renderUnit := range renderUnits {
		if (renderUnit.Status == "published") {
			result = append(result, renderUnit)
		}
	}

	return result, nil
}

func (r *Repository) GetDraftRenderUnit() (RenderUnit, error) {
	renderUnits, err := r.GetRenderUnits()
	if err != nil {
		return RenderUnit{}, err
	}

	var result RenderUnit
	for _, renderUnit := range renderUnits {
		if (renderUnit.Status == "draft") {
			result = renderUnit
		}
	}

	return result, nil
}


func (r *Repository) GetNextPublishedRenderUnitTo(id int) (RenderUnit, error) {
	renderUnits, err := r.GetPublishedRenderUnits()
	if err != nil {
		return RenderUnit{}, err
	}

	var currentIndex int
	for i, renderUnit := range renderUnits {
		if (renderUnit.ID == id) {
			currentIndex = i
		}
	}

	if (currentIndex + 1) == len(renderUnits) {
		return renderUnits[currentIndex], nil
	}

	var result RenderUnit = renderUnits[currentIndex + 1]

	return result, nil
}