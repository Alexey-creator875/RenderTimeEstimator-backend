package repository

import (
	"RenderTimeEstimator/internal/app/ds"
)

func (r *Repository) AddLike(like *ds.Likes) error {
	err := r.db.Create(like).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteLike(userID uint, renderServerUnitID uint) error {
	err := r.db.
		Where("user_id = ? AND render_server_unit_id = ?", userID, renderServerUnitID).
		Delete(&ds.Likes{}).Error

	if err != nil {
		return err
	}

	return nil
}
