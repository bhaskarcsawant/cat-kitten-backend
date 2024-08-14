package controllers

import (
	"encoding/json"
	"net/http"
	models "server/model"
	services "server/service"

	mux "github.com/gorilla/mux"
)

// SetUserPointsHandler handles setting user points
func SetUserPointsHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserScore
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addedData, err := services.SetUserGamePoints(user.Username, user.Points)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(addedData)
}

// IncrementUserPointsHandler handles incrementing user points
func IncrementUserPointsHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserScore

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addedData, err := services.IncrementUserGamePoints(user.Username, int(user.Points))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(addedData)
}

// GetAllUserPointsHandler retrieves all user points in descending order
func GetAllUserPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["status"]
	users := services.GetAllUserPointsDesc(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
