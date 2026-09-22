package repository

import (
	"RenderTimeEstimator/internal/app/ds"
	"errors"
	"fmt"

	"gorm.io/gorm"
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

func (r *Repository) GetDraftRenderServerUnit(creatorID uint) (ds.RenderServerUnit, bool, error) {
	draftRenderServerUnit := ds.RenderServerUnit{}
	err := r.db.Where("status = ? AND creator_id = ?", "draft", creatorID).Take(&draftRenderServerUnit).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.RenderServerUnit{}, false, nil
		} else {
			return ds.RenderServerUnit{}, false, err
		}
	}

	return draftRenderServerUnit, true, nil
}


func (r *Repository) GetNextPublishedRenderServerUnitTo(id int) (ds.RenderServerUnit, error) {
	var next ds.RenderServerUnit
	err := r.db.Where("id > ?", id).Order("id ASC").First(&next).Error

	if err != nil {
		return ds.RenderServerUnit{}, err
	}

	return next, nil
}

func (r *Repository) AddDraftRenderServerUnit(renderServerUnit ds.RenderServerUnit) error {
	err := r.db.Create(&renderServerUnit).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateRenderServerUnit(renderServerUnit ds.RenderServerUnit) error {
	err := r.db.Save(&renderServerUnit).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteRenderServerUnit(id int) error {
	query := "UPDATE render_server_units SET status = $1 WHERE id = $2"

    err := r.db.Exec(query, "deleted", id).Error

    if err != nil {
        return err
    }

    return nil
}

func (r *Repository) GetLikesNumber(id int) (int, error) {
    query := "SELECT COUNT(*) FROM likes WHERE render_server_unit_id = $1"

    row := r.db.Raw(query, id).Row()

    var count int
	err := row.Scan(&count)

    if  err != nil {
        return 0, err
    }

    return count, nil
}
