// main.go in user-service directory
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Set up the router.
	r := mux.NewRouter()

	// Start the server.
	log.Println("Starting user service on port 5000...")
	if err := http.ListenAndServe(":5000", r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
