package main

import (
	"fmt"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {

	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	// mux.Handle("/app/", http.FileServer(http.Dir(filepathRoot)))

	mux.HandleFunc("/healthz", healthHandler)

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
