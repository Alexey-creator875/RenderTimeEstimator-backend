package ds

import "time"

type RenderServerUnit struct {
	ID          int			`gorm:"primaryKey" json:"id"`
	Processor   string		`gorm:"type:varchar(30)"  json:"processor"`
	Description string		`gorm:"type:varchar(300)" json:"description"`
	Status      string		`gorm:"type:varchar(10)"  json:"status"`
	Image       string		`gorm:"type:varchar(100)" json:"image"`
	Video       string		`gorm:"type:varchar(100)" json:"video"`
	Cores       int			`json:"cores"`
	RAM         int			`json:"ram"`
							
	CreatedAt	time.Time	`gorm:"autoCreateTime" json:"-"`
	CreatorID	uint		`gorm:"index"           json:"-"`
	Creator		User		`gorm:"foreignKey:CreatorID" json:"-"`
	FormedAt	time.Time	`gorm:"autoUpdateTime"  json:"-"`
}
