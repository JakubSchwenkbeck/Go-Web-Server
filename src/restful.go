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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	defer mu.Unlock()

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users[user.ID] = user
	w.WriteHeader(http.StatusCreated)
	defer mu.Unlock()

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()

	id := mux.Vars(r)["id"]
	user, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	defer mu.Unlock()

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()

	id := mux.Vars(r)["id"]
	if _, exists := users[id]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
	defer mu.Unlock()

}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/users", CreateUser).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", GetUser).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", DeleteUser).Methods("DELETE")
}
