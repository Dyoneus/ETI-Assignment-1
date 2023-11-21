// By Ong Jia Yuan S10227735B
// main.go in trip-service directory
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
	// Make sure to change the port number if both services will run on the same host.
	log.Println("Starting trip service on port 5001...")
	if err := http.ListenAndServe(":5001", r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
