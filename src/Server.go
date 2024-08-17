package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	// define home page

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
	// define header function

	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func RequestInformation(w http.ResponseWriter, r *http.Request) {
	for name, headers := range r.Header {
		for _, h := range headers {
			// get header informations and print them into readable format for information about client / request

			// get platform info
			if strings.Contains(name, "Sec-Ch-Ua-Platform") {
				h = strings.Trim(h, "")
				fmt.Fprintf(w, "Clinet running on a "+h+" platform \n")
			}
			if strings.Contains(name, "Sec-Ch-Ua-Mobile") {
				if h != "?0" {
					fmt.Fprintf(w, "Client IS running on a mobile device \n")
				} else {
					fmt.Fprintf(w, "Client NOT running on a mobile device \n")
				}
			}

		}
	}

}

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/header", header)
	http.HandleFunc("/ClientInfo", RequestInformation)

	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
