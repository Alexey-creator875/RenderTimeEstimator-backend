package main

import (
	"RenderTimeEstimator/internal/app/ds"
	"RenderTimeEstimator/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.RenderServerUnit{},
		&ds.User{},
		&ds.Likes{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}