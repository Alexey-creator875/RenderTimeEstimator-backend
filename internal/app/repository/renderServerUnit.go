package repository

import (
	"RenderTimeEstimator/internal/app/ds"
	"fmt"
)

func (r *Repository) GetRenderServerUnits() ([]ds.RenderServerUnit, error) {
	var renderServerUnits []ds.RenderServerUnit
	err := r.db.Find(&renderServerUnits).Error

	if err != nil {
		return nil, err
	}
	if len(renderServerUnits) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return renderServerUnits, nil
}

func (r *Repository) GetRenderServerUnit(id int) (ds.RenderServerUnit, error) {
	renderServerUnit := ds.RenderServerUnit{}
	err := r.db.Where("id = ?", id).First(&renderServerUnit).Error

	if err != nil {
		return ds.RenderServerUnit{}, err
	}

	return renderServerUnit, nil
}

func (r *Repository) GetPublishedRenderServerUnits() ([]ds.RenderServerUnit, error) {
	var renderServerUnits []ds.RenderServerUnit
	err := r.db.Where("status = ?", "published").Find(&renderServerUnits).Error

	if err != nil {
		return nil, err
	}

	return renderServerUnits, nil
}

func (r *Repository) GetPublishedRenderServerUnitsByRAM(min_ram int, max_ram int) ([]ds.RenderServerUnit, error) {
	var renderServerUnits []ds.RenderServerUnit
	err := r.db.Where("status = ? AND RAM BETWEEN ? AND ?", "published", min_ram, max_ram).Find(&renderServerUnits).Error

	if err != nil {
		return nil, err
	}

	return renderServerUnits, nil
}

func (r *Repository) GetDraftRenderServerUnit() (ds.RenderServerUnit, error) {
	renderServerUnit := ds.RenderServerUnit{}
	err := r.db.Where("status = ?", "draft").First(&renderServerUnit).Error

	if err != nil {
		return ds.RenderServerUnit{}, err
	}

	return renderServerUnit, nil
}


func (r *Repository) GetNextPublishedRenderServerUnitTo(id int) (ds.RenderServerUnit, error) {
	var next ds.RenderServerUnit
	err := r.db.Where("id > ?", id).Order("id ASC").First(&next).Error

	if err != nil {
		return ds.RenderServerUnit{}, err
	}

	return next, nil
}