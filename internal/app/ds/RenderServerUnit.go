package ds

type RenderServerUnit struct {
  	ID			int		`gorm:"primaryKey"`
	Status		string	`gorm:"type:varchar(10)"`
	Processor 	string	`gorm:"type:varchar(30)"`
	Cores 		int
	RAM 		int
	Description string	`gorm:"type:varchar(300)"`
	Image	 	string	`gorm:"type:varchar(100)"`
	Video	 	string	`gorm:"type:varchar(100)"`
}
