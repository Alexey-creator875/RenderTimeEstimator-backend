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

func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	user := &ds.User{
		Login: "login",
	}

	err := r.db.First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}