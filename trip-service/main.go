// By Ong Jia Yuan S10227735B
// /trip-service/main.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the router.
	r := mux.NewRouter()

	// The setup code would be very similar to the user-service's main.go
	log.Println("Starting trip service on port 5001...")
	if err := http.ListenAndServe(":5001", r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
