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
	fmt.Fprintf(w, "Hello, World!\n")
	fmt.Fprintf(w, "This page has been accessed "+strconv.Itoa(currentCount)+" times!")
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
	form := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Register</title>
	</head>
	<body>
		<h1>Register</h1>
		<form action="/register" method="post">
			<label for="id">ID:</label>
			<input type="text" id="id" name="id" required><br><br>
			<label for="name">Name:</label>
			<input type="text" id="name" name="name" required><br><br>
			<label for="password">Password:</label>
			<input type="password" id="password" name="password" required><br><br>
			<input type="submit" value="Register">
		</form>
	</body>
	</html>`
	RenderHTML(w, r, form)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
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
}

// LoginPage serves the user login page
func LoginPage(w http.ResponseWriter, r *http.Request) {
	form := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Login</title>
	</head>
	<body>
		<h1>Login</h1>
		<form action="/login" method="post">
			<label for="username">Username:</label>
			<input type="text" id="username" name="username" required><br><br>
			<label for="password">Password:</label>
			<input type="password" id="password" name="password" required><br><br>
			<input type="submit" value="Login">
		</form>
	</body>
	</html>`
	RenderHTML(w, r, form)
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
		<form action="/send" method="post">
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
}
