// By Ong Jia Yuan S10227735B
// main.go in user-service directory
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Initialize the database connection.
	dsn := "username:password@tcp(127.0.0.1:5000)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Make sure to defer the closing of the database connection.
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Could not get database: %v", err)
	}
	defer sqlDB.Close()

	// Set up the router.
	r := mux.NewRouter()

	// Start the server.
	log.Println("Starting user service on port 5000...")
	if err := http.ListenAndServe(":5000", r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
