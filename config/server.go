package server

import (
	"fmt"
	"net/http"
	controllers "server/controller"
	services "server/service"
)

// CORS middleware to handle CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust the origin as needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			return
		}

		// Proceed with the actual request
		next.ServeHTTP(w, r)
	})
}

// StartServer initializes the server and routes
func StartServer() error {
	// Initialize Redis client
	services.InitializeRedisClient()

	// HTTP handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/set", controllers.SetUserPointsHandler)
	mux.HandleFunc("/increment", controllers.IncrementUserPointsHandler)
	mux.HandleFunc("/getall", controllers.GetAllUserPointsHandler)

	// Apply CORS middleware
	handler := CORSMiddleware(mux)

	// Start the server
	fmt.Println("Server is running on port 8080...")
	return http.ListenAndServe(":8080", handler)
}
