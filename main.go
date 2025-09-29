package main

import (
	"fmt"
	"net/http"
)

func main() {

	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Starting server on :8080")
	err := srv.ListenAndServe()

	if err != nil {
		fmt.Printf("couldn't list and serve: %v", err)
		return
	}
}
