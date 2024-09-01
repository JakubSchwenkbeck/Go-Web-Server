package main

/** IMPORTS */
import (
	"fmt"
	"log"
	"net/http"
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

// PAGES ARE DECLARED IN pages.go

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

	// Write enhanced HTML content with inline CSS for better styling
	fmt.Fprintf(w, `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Server Information</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f9f9f9;
                color: #333;
                margin: 0;
                padding: 0;
                display: flex;
                flex-direction: column;
                align-items: center;
                text-align: center;
            }
            .container {
                max-width: 800px;
                margin: 20px;
                padding: 20px;
                background: #fff;
                border-radius: 8px;
                box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            }
            h1 {
                color: #007BFF;
                margin-bottom: 20px;
            }
            p {
                line-height: 1.6;
                margin: 10px 0;
            }
            a {
                color: #007BFF;
                text-decoration: none;
            }
            a:hover {
                text-decoration: underline;
            }
            .info-box {
                background: #f1f1f1;
                border-left: 5px solid #007BFF;
                padding: 15px;
                margin: 20px 0;
                border-radius: 5px;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>General Information about this Server</h1>
            <div class="info-box">
                <p><strong>Programmed and run by:</strong> Jakub</p>
                <p><strong>Focus:</strong> Go Lang for Backend and Web Server development</p>
                <p><strong>GitHub Profile:</strong> <a href="https://github.com/jakubschwenkbeck" target="_blank">Jakub's GitHub</a></p>
            </div>
            <div class="info-box">
                <p><strong>Technologies Used:</strong></p>
                <ul>
                    <li>Go Language</li>
                    <li>MySQL</li>
                    <li>JWT Authentication</li>
                    <li>RESTful API</li>
                    <li>HTML/CSS</li>
                </ul>
            </div>
            
        </div>
    </body>
    </html>`)
}

func main() {
	r := mux.NewRouter()

	// Register routes
	r.HandleFunc("/", HomePage)
	r.HandleFunc("/header", header)
	r.HandleFunc("/ClientInfo", RequestInformation)
	r.HandleFunc("/info", GeneralInformation)
	r.HandleFunc("/register", RegisterPage)
	r.HandleFunc("/restful/register", RegisterUser).Methods("POST")
	r.HandleFunc("/login", LoginPage)
	r.HandleFunc("/restful/login", LoginUser).Methods("POST")
	r.HandleFunc("/send", SendMessagePage)
	r.HandleFunc("/restful/send", SendMessage).Methods("POST")

	// Register routes from restful.go
	RegisterRoutes(r)

	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)

	connectDB()

	ChatAppMain(*r, port)

	log.Fatal(http.ListenAndServe(":"+port, r)) // Use the router here

}
