package controllers

import (
	services "ProjectHash/services"
	"encoding/json"
	"net/http"
)

type validatePasswordRequest struct {
	HashedPassword string `json:"hashed_password"`
	Password       string `json:"password"`
}

type validatePasswordResponse struct {
	Valid bool   `json:"valid"`
	Err   string `json:"error,omitempty"`
}

func ValidatePasswordHandler(svc services.PasswordService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request validatePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		valid, err := svc.ValidatePassword(request.HashedPassword, request.Password)
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		response := validatePasswordResponse{
			Valid: valid,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
