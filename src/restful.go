package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	// Map to store users with their ID as the key
	users = make(map[string]User)
)

/**
 * Implementing CRUD (Create, Read, Update, Delete) operations for user management.
 * This simulates a RESTful API for handling users.
 */

/**
 * GetUsers handles the HTTP GET request for retrieving all users.
 * It locks the mutex to ensure thread safety while accessing the shared `users` map.
 * Responds with a JSON-encoded list of all users.
 *
 * @param w http.ResponseWriter: The response writer to send data to the client.
 * @param r *http.Request: The incoming HTTP request.
 */
func GetUsers(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

/**
 * CreateUser handles the HTTP POST request for creating a new user.
 * It locks the mutex to ensure thread safety while modifying the `users` map.
 * The user data is decoded from the request body and added to the `users` map.
 * Responds with a status code 201 Created if successful.
 *
 * @param w http.ResponseWriter: The response writer to send data to the client.
 * @param r *http.Request: The incoming HTTP request.
 */
func CreateUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users[user.ID] = user
	defer mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}

/**
 * GetUser handles the HTTP GET request for retrieving a specific user by ID.
 * It locks the mutex to ensure thread safety while accessing the `users` map.
 * Responds with a JSON-encoded user object if found, or a 404 Not Found status if the user does not exist.
 *
 * @param w http.ResponseWriter: The response writer to send data to the client.
 * @param r *http.Request: The incoming HTTP request.
 */
func GetUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	id := mux.Vars(r)["id"]
	user, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

/**
 * UpdateUser handles the HTTP PUT request for updating a specific user's information.
 * It locks the mutex to ensure thread safety while modifying the `users` map.
 * The updated user data is decoded from the request body and the existing user is updated.
 * Responds with the updated user object if successful, or a 404 Not Found status if the user does not exist.
 *
 * @param w http.ResponseWriter: The response writer to send data to the client.
 * @param r *http.Request: The incoming HTTP request.
 */
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()

	id := mux.Vars(r)["id"]
	user, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the user's name in the map
	user.Name = updatedUser.Name
	users[id] = user
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

/**
 * DeleteUser handles the HTTP DELETE request for deleting a specific user by ID.
 * It locks the mutex to ensure thread safety while modifying the `users` map.
 * The user is removed from the `users` map.
 * Responds with a status code 204 No Content if successful, or a 404 Not Found status if the user does not exist.
 *
 * @param w http.ResponseWriter: The response writer to send data to the client.
 * @param r *http.Request: The incoming HTTP request.
 */
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()

	id := mux.Vars(r)["id"]
	if _, exists := users[id]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	delete(users, id)
	defer mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

/**
 * RegisterRoutes sets up the routing for the user management API.
 * It registers the CRUD operations with specific HTTP methods and paths.
 *
 * @param r *mux.Router: The router to which routes are added.
 */
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", GetUsers).Methods("GET")            // Route for listing all users (accessible by all)
	r.HandleFunc("/users/{id:[0-9]+}", GetUser).Methods("GET") // Route for retrieving a user by ID (accessible by all)

	r.Handle("/users", IsAdmin(http.HandlerFunc(CreateUser))).Methods("POST")               // Route for creating a new user (admin only)
	r.Handle("/users/{id:[0-9]+}", IsAdmin(http.HandlerFunc(DeleteUser))).Methods("DELETE") // Route for deleting a user by ID (admin only)
}
