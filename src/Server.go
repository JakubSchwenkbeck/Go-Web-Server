package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/", helloHandler)
	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
