package repository

import (
	"fmt"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type RenderServerUnit struct {
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

func (r *Repository) GetRenderServerUnits() ([]RenderServerUnit, error) {
	renderServerUnits := []RenderServerUnit{
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

	if len(renderServerUnits) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return renderServerUnits, nil
}

func (r *Repository) GetRenderServerUnit(id int) (RenderServerUnit, error) {
	renderServerUnits, err := r.GetRenderServerUnits()
	if err != nil {
		return RenderServerUnit{}, err
	}

	for _, renderServerUnit := range renderServerUnits {
		if renderServerUnit.ID == id {
			return renderServerUnit, nil
		}
	}
	return RenderServerUnit{}, fmt.Errorf("заказ не найден")
}

func (r *Repository) GetPublishedRenderServerUnitsByRAM(min_ram int, max_ram int) ([]RenderServerUnit, error) {
	renderServerUnits, err := r.GetRenderServerUnits()
	if err != nil {
		return []RenderServerUnit{}, err
	}

	var result []RenderServerUnit
	for _, renderServerUnit := range renderServerUnits {
		if renderServerUnit.RAM >= min_ram && renderServerUnit.RAM <= max_ram && renderServerUnit.Status == "published" {
			result = append(result, renderServerUnit)
		}
	}

	return result, nil
}

func (r *Repository) GetPublishedRenderServerUnits() ([]RenderServerUnit, error) {
	renderServerUnits, err := r.GetRenderServerUnits()
	if err != nil {
		return []RenderServerUnit{}, err
	}

	var result []RenderServerUnit
	for _, renderServerUnit := range renderServerUnits {
		if (renderServerUnit.Status == "published") {
			result = append(result, renderServerUnit)
		}
	}

	return result, nil
}

func (r *Repository) GetDraftRenderServerUnit() (RenderServerUnit, error) {
	renderServerUnits, err := r.GetRenderServerUnits()
	if err != nil {
		return RenderServerUnit{}, err
	}

	var result RenderServerUnit
	for _, renderServerUnit := range renderServerUnits {
		if (renderServerUnit.Status == "draft") {
			result = renderServerUnit
		}
	}

	return result, nil
}


func (r *Repository) GetNextPublishedRenderServerUnitTo(id int) (RenderServerUnit, error) {
	renderServerUnits, err := r.GetPublishedRenderServerUnits()
	if err != nil {
		return RenderServerUnit{}, err
	}

	var currentIndex int
	for i, renderServerUnit := range renderServerUnits {
		if (renderServerUnit.ID == id) {
			currentIndex = i
		}
	}

	if (currentIndex + 1) == len(renderServerUnits) {
		return renderServerUnits[currentIndex], nil
	}

	var result RenderServerUnit = renderServerUnits[currentIndex + 1]

	return result, nil
}