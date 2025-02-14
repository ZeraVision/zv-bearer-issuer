package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Target")
}

func StoreData(w http.ResponseWriter, r *http.Request) {
	// Set cors headers
	setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	//* Can do authentication here as required for your use case *//

	process(w, r)
}

// StartHTTPServer starts the HTTP server to listen for data on a specified port
func StartHTTPServer(port string) {
	http.HandleFunc("/store", StoreData)

	fmt.Println("HTTP Server started. Listening on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func StartAPI() {
	port := os.Getenv("API_PORT")
	if len(port) < 1 {
		port = "80"
	}

	go StartHTTPServer(port)
}
