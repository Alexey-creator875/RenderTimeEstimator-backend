package repository

import (
	"RenderTimeEstimator/internal/app/ds"
	"github.com/google/uuid"
)


func (r *Repository) AddUser(user *ds.User) error {
	err := r.db.Create(user).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Register(user *ds.User) error {
	if user.UUID == uuid.Nil {
		user.UUID = uuid.New()
	}

	return r.db.Create(user).Error
}