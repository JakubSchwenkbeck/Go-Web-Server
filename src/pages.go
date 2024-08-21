package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HomePage serves the landing page of the server
func HomePage(w http.ResponseWriter, r *http.Request) {
	semaphore <- struct{}{}
	defer func() { <-semaphore }() // Release semaphore slot

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	counter++
	currentCount := counter
	mu.Unlock()

	// HTML content
	htmlContent := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Welcome to My Homepage</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f0f8ff;
                color: #333;
                text-align: center;
                padding: 50px;
            }
            h1 {
                color: #4CAF50;
            }
            .counter {
                font-size: 24px;
                margin: 20px 0;
            }
            .button {
                background-color: #4CAF50;
                border: none;
                color: white;
                padding: 15px 32px;
                text-align: center;
                text-decoration: none;
                display: inline-block;
                font-size: 16px;
                margin: 4px 2px;
                cursor: pointer;
                border-radius: 4px;
                transition: background-color 0.3s;
            }
            .button:hover {
                background-color: #45a049;
            }
            img.logo {
                width: 300px;
                border-radius: 8px;
            }
        </style>
    </head>
    <body>
        <h1>Welcome to My Homepage!</h1>
        <p class="counter">This page has been accessed <strong>` + strconv.Itoa(currentCount) + `</strong> times!</p>
        <p>
			 <div>
			 <a href="/" class="button">Refresh Page</a>
                <a href="/register" class="button">Register</a>
                <a href="/login" class="button">Login</a>
				<a href="/info" class="button">Information</a>

            </div>
        </p>
		<img src="https://i.imgur.com/p0pVfQB.jpg" alt="Logo" class="logo">
    </body>
    </html>
    `

	// Send the HTML response
	fmt.Fprintf(w, htmlContent)
}

// RenderHTML serves a basic HTML form page
func RenderHTML(w http.ResponseWriter, r *http.Request, tmpl string) {
	t, err := template.New("form").Parse(tmpl)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

// RegisterPage serves the user registration page
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		form := `
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Register</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: white;
                    padding: 30px;
                    border-radius: 10px;
                    box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
                    width: 100%;
                    max-width: 400px;
                }
                h1 {
                    text-align: center;
                    color: #333;
                }
                label {
                    font-weight: bold;
                    color: #333;
                    display: block;
                    margin-top: 10px;
                }
						img.logo {
           			 width: 320px; /* Adjusted size */
           			 margin-bottom: 20px;
       				 }
     	   	
                input[type="text"], input[type="password"] {
                    width: 100%;
                    padding: 10px;
                    margin: 10px 0;
                    border: 1px solid #ccc;
                    border-radius: 5px;
                    box-sizing: border-box;
                }
                input[type="submit"] {
                    width: 100%;
                    padding: 10px;
                    background-color: #007BFF;
                    color: white;
                    border: none;
                    border-radius: 5px;
                    cursor: pointer;
                    font-size: 16px;
                    margin-top: 20px;
                }
                input[type="submit"]:hover {
                    background-color: #0056b3;
                }
            </style>
        </head>
        <body>
            <div class="container">
				<img src="https://i.imgur.com/p0pVfQB.jpg" alt="Logo" class="logo">

                <h1>Register</h1>
                <form action="/restful/register" method="post">
                    <label for="id">ID:</label>
                    <input type="text" id="id" name="id" required>

                    <label for="name">Name:</label>
                    <input type="text" id="name" name="name" required>

                    <label for="password">Password:</label>
                    <input type="password" id="password" name="password" required>

                    <input type="submit" value="Register">
                </form>
            </div>
        </body>
        </html>`
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, form)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// RegisterUser handles form submission for user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		name := r.FormValue("name")
		password, err := HashPassword(r.FormValue("password"))
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users (id, name, password) VALUES (?, ?, ?)", id, name, password)
		if err != nil {
			http.Error(w, "Error registering user", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User %s registered successfully!", name)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}

// Login Page with HTML
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		form := `
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Login</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: white;
                    padding: 30px;
                    border-radius: 10px;
                    box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
                    width: 100%;
                    max-width: 400px;
                }
                h1 {
                    text-align: center;
                    color: #333;
                }
                label {
                    font-weight: bold;
                    color: #333;
                    display: block;
                    margin-top: 10px;
                }
					img.logo {
            width: 320px; /* Adjusted size */
            margin-bottom: 20px;
        }
        form {
            display: flex;
            flex-direction: column;
            align-items: center;
        }
                input[type="text"], input[type="password"] {
                    width: 100%;
                    padding: 10px;
                    margin: 10px 0;
                    border: 1px solid #ccc;
                    border-radius: 5px;
                    box-sizing: border-box;
                }
                input[type="submit"] {
                    width: 100%;
                    padding: 10px;
                    background-color: #007BFF;
                    color: white;
                    border: none;
                    border-radius: 5px;
                    cursor: pointer;
                    font-size: 16px;
                    margin-top: 20px;
                }
                input[type="submit"]:hover {
                    background-color: #0056b3;
                }
            </style>
        </head>
        <body>
            <div class="container">
			<img src="https://i.imgur.com/p0pVfQB.jpg" alt="Logo" class="logo">

                <h1>Login</h1>
                <form action="/restful/login" method="post">
                    <label for="username">Username:</label>
                    <input type="text" id="username" name="username" required>

                    <label for="password">Password:</label>
                    <input type="password" id="password" name="password" required>

                    <input type="submit" value="Login">
                </form>
            </div>
        </body>
        </html>`
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, form)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// LoginUser handles user login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE name = ?", username).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, _ := HashPassword(password)
	fmt.Fprintf(w, "Login successful! Your token is: %s", token)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// SendMessagePage serves the send message form
func SendMessagePage(w http.ResponseWriter, r *http.Request) {
	form := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Send Message</title>
	</head>
	<body>
		<h1>Send Message</h1>
		<form action="/restful/send" method="post">
			<label for="senderID">Sender ID:</label>
			<input type="text" id="senderID" name="senderID" required><br><br>
			<label for="receiverID">Receiver ID:</label>
			<input type="text" id="receiverID" name="receiverID" required><br><br>
			<label for="message">Message:</label>
			<textarea id="message" name="message" required></textarea><br><br>
			<input type="submit" value="Send Message">
		</form>
	</body>
	</html>`
	RenderHTML(w, r, form)
}

// SendMessage handles sending a message
func SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	senderID := r.FormValue("senderID")
	receiverID := r.FormValue("receiverID")
	message := r.FormValue("message")

	_, err := db.Exec("INSERT INTO messages (sender_id, receiver_id, message, timestamp) VALUES (?, ?, ?, ?)",
		senderID, receiverID, message, time.Now())
	if err != nil {
		http.Error(w, "Error sending message", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Message sent from %s to %s", senderID, receiverID)
	http.Redirect(w, r, "/send", http.StatusSeeOther)
}
