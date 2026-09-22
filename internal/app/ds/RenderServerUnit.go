package ds

import "time"

type RenderServerUnit struct {
	ID          int       `gorm:"primaryKey"`
	Processor   string    `gorm:"type:varchar(30)"`
	Description string    `gorm:"type:varchar(300)"`
	Status      string    `gorm:"type:varchar(10)"`
	Image       string    `gorm:"type:varchar(100)"`
	Video       string    `gorm:"type:varchar(100)"`
	Cores       int
	RAM         int

	CreatedAt   time.Time `gorm:"autoCreateTime"`
	CreatorID   uint       `gorm:"index"`
	Creator     User      `gorm:"foreignKey:CreatorID"`
	FormedAt    time.Time `gorm:"autoUpdateTime"`
}
