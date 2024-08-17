package main

/** IMPORTS */
import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

var ( // shared resource
	counter int
	mu      sync.Mutex
	// Semaphore with a capacity of 5 for accesses
	semaphore = make(chan struct{}, 5)
)

/** HomePage creates the landing page which is first seen when logging onto server  */

func HomePage(w http.ResponseWriter, r *http.Request) {
	// define home page

	// Acquire semaphore slot
	semaphore <- struct{}{}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Lock and update shared resource
	mu.Lock()
	counter++
	currentCount := counter
	mu.Unlock()
	fmt.Fprintf(w, "Hello, World!\n")
	fmt.Fprintf(w, "This page has been accessed "+strconv.Itoa(currentCount)+" times!")

	defer func() { <-semaphore }() // Release semaphore slot

}

/** header is more for information about the client/request, might delete later */
func header(w http.ResponseWriter, r *http.Request) {
	// define header function

	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

/** Use parts of the header information and print in readable format*/
func RequestInformation(w http.ResponseWriter, r *http.Request) {
	for name, headers := range r.Header {
		for _, h := range headers {
			// get header informations and print them into readable format for information about client / request

			// get platform info
			if strings.Contains(name, "Sec-Ch-Ua-Platform") {
				h = strings.Trim(h, "")
				fmt.Fprintf(w, "Client running on a "+h+" platform \n")
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

/** Generate an HTML site for the General Information  */

func GeneralInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Write HTML content with dynamic information
	fmt.Fprintf(w, "<!DOCTYPE html>")
	fmt.Fprintf(w, "<html lang=\"en\">")
	fmt.Fprintf(w, "<head>")
	fmt.Fprintf(w, "<meta charset=\"UTF-8\">")
	fmt.Fprintf(w, "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">")
	fmt.Fprintf(w, "<title>Server Information</title>")
	fmt.Fprintf(w, "</head>")
	fmt.Fprintf(w, "<body>")
	fmt.Fprintf(w, "<h1>General Information about this Server</h1>")
	fmt.Fprintf(w, "<p>Programmed and run by Jakub</p>")
	fmt.Fprintf(w, "<p>Focuses on Go Lang for Backend and Web Server development</p>")
	fmt.Fprintf(w, "<p>Github: <a href=\"https://github.com/jakubschwenkbeck\">Jakub's GitHub</a></p>")
	fmt.Fprintf(w, "</body>")
	fmt.Fprintf(w, "</html>")
}

func main() {
	r := mux.NewRouter()

	// Register routes
	r.HandleFunc("/", HomePage)
	r.HandleFunc("/header", header)
	r.HandleFunc("/ClientInfo", RequestInformation)
	r.HandleFunc("/info", GeneralInformation)

	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r)) // Use the router here
}
