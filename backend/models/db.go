package models

import (
	"log"

	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var validate = validator.New()
var gormDB *gorm.DB = InitDb()

func InitDb() *gorm.DB {
	dsn := "user=postgres password=dbpassword dbname=mood host=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("db error:", err)
	}
	db.AutoMigrate(
		&User{},
		&AuthToken{},
		&MonAgent{},
		&AnalysedPage{},
	)

	return db
}
