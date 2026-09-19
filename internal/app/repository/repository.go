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
	Likes []int
	ImagePath string
	VideoPath string
}

func (r *Repository) GetRenderServerUnits() ([]RenderServerUnit, error) {
	renderServerUnits := []RenderServerUnit{
		{
			ID:    1,
			Status: "published",
			Processor: "RTX 2080 Ti",
			Cores: 20,
			RAM: 16,
			Description: "Надёжный выбор для рендеринга и игр. Отличная производительность в 2K и 4K, поддерживает все современные технологии трассировки лучей. Идеально для студийных проектов средней сложности.",
			Likes: []int{12, 45, 78, 23, 89, 34, 56, 91, 67, 10, 42, 88, 15, 73, 29, 61, 95, 37, 54, 82},
			ImagePath: "2080-Ti.jpg",
			VideoPath: "2080-Ti.mp4",
		},
		{
			ID:    2,
			Status: "published",
			Processor: "RTX 3090",
			Cores: 20,
			RAM: 24,
			Description: "Флагманская карта с 24 ГБ видеопамяти. Мгновенный рендеринг сложных сцен, работа с нейросетями и 8K-видео. Профессиональный инструмент для 3D-художников и дата-сайентистов.",
			Likes: []int{5, 18, 33, 47, 62, 71, 84, 99, 11, 26, 39, 52, 68, 75, 83, 94, 7, 21, 36, 49, 64, 77, 86, 92, 14},
			ImagePath: "3090.jpg",
			VideoPath: "3090.mp4",
		},
		{
			ID:    3,
			Status: "published",
			Processor: "RTX 4090",
			Cores: 20,
			RAM: 24,
			Description: "Абсолютный лидер производительности. Справляется с самыми тяжёлыми задачами: симуляция физики, объёмная визуализация, рендеринг в реальном времени. Будущее уже сегодня.",
			Likes: []int{101, 115, 123, 134, 142, 156, 167, 178, 189, 195, 104, 117, 126, 138, 149, 158, 169, 176, 187, 198, 109, 112, 128, 136, 147, 154, 165, 172, 183, 192},
			ImagePath: "4090.png",
			VideoPath: "4090.mp4",
		},
		{
			ID:    4,
			Status: "draft",
			Processor: "",
			Cores: 0,
			RAM: 0,
			Description: "Описание",
			Likes: []int{},
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
			Likes: []int{2, 9, 16, 23, 30, 37, 44, 51, 58, 65, 72, 79, 86, 93, 100, 107, 114, 121, 128, 135, 142, 149, 156, 163, 170, 177, 184, 191, 198, 205, 212, 219, 226, 233, 240, 247, 254, 261, 268, 275},
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