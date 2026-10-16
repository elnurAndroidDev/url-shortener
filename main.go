package main

import (
	"fmt"
	"net/http"
	"url-shortener/handlers"
	"url-shortener/storage"
)

func main() {
	store := storage.NewStore()
	handler := handlers.NewHandler(store)

	http.HandleFunc("/shorten", handler.ShortenURL)
	http.HandleFunc("/",handler.RedirectURL)

	port := ":8080"
	fmt.Println("Server is running on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
