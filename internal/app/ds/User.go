package ds

import (
	"RenderTimeEstimator/internal/app/role"

	"github.com/google/uuid"
)

type User struct {
	ID       	uint		`gorm:"primaryKey"`
	UUID		uuid.UUID	`gorm:"type:uuid"`
	Login		string		`gorm:"type:varchar(40)" json:"login"`
	Password 	string		`gorm:"type:varchar(40)" json:"-"`
	Role 		role.Role	`sql:"type:string;"`
}

