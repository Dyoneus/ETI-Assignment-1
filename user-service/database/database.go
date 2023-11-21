// user-service/database/database.go
package database

import (
	"user-service/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitializeDatabase sets up the database connection and runs migrations.
func InitializeDatabase() (*gorm.DB, error) {
	dsn := "username:password@tcp(127.0.0.1:3306)/carpool?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migrations
	db.AutoMigrate(&models.User{}, &models.CarOwnerProfile{})

	return db, nil
}
