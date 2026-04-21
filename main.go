package main

import (
	"ascii-art-web-export-file/internal/handlers"
	"log"
	"net/http"
)

func main() {
	// main only bootstraps the server and delegates route setup to the
	// handlers package so HTTP logic stays out of the entrypoint.
	//link css with go
	mux := http.NewServeMux()
	handlers.Register(mux)
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("image"))))

	log.Println("Server started at http://localhost:8080")
	// The mux is the router that sends each request to the matching handler.
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
