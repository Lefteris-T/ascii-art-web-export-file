package main

import (
	"ascii-art-web-export-file/internal/handlers"
	"log"
	"net/http"
)

// main bootstraps the web server, registers application routes, and exposes
// static assets used by the templates.
func main() {
	mux := http.NewServeMux()
	handlers.Register(mux)

	// Serve static files separately from the application handlers.
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("image"))))

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
