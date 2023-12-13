// trip-service/database/database.go
package database

import (
	"trip-service/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeDatabase() (*gorm.DB, error) {
	dsn := "user:password@tcp(127.0.0.1:3306)/carpool_trips?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migrations
	db.AutoMigrate(&models.Trip{}, &models.Reservation{})

	return db, nil
}
