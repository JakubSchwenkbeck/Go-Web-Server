package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
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

func header(w http.ResponseWriter, r *http.Request) {
	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func platform(w http.ResponseWriter, r *http.Request) {
	for name, headers := range r.Header {
		for _, h := range headers {

			if strings.Contains(name, "Sec-Ch-Ua-Platform") {
				fmt.Fprintf(w, "%v: %v\n", "Platform: ", h)
			}
		}
	}

}

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/header", header)
	http.HandleFunc("/platform", platform)

	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
