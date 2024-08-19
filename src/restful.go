package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Defining JSON User strcuture with id and name
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	users = make(map[string]User) // map for users

)

/** Implementing CRUD (Create, Read, Update, Delete) while simulating the restful API for users*/

func GetUsers(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)

}

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

	// Update the user in the map
	user.Name = updatedUser.Name
	users[id] = user
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}

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
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", GetUsers).Methods("GET")            // accessible by all
	r.HandleFunc("/users/{id:[0-9]+}", GetUser).Methods("GET") // accessible by all

	r.Handle("/users", IsAdmin(http.HandlerFunc(CreateUser))).Methods("POST")               // admin only
	r.Handle("/users/{id:[0-9]+}", IsAdmin(http.HandlerFunc(DeleteUser))).Methods("DELETE") // admin only
}
