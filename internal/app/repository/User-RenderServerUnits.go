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

func (r *Repository) IsLikedByUser(renderServerUnitId int, userId int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM likes WHERE user_id = $1 AND render_server_unit_id = $2)"
	row := r.db.Raw(query, userId, renderServerUnitId).Row()

	var exists bool
	err := row.Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) GetLikesNumber(id int) (int, error) {
	query := "SELECT COUNT(*) FROM likes WHERE render_server_unit_id = $1"

	row := r.db.Raw(query, id).Row()

	var count int
	err := row.Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}