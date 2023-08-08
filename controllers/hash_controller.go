package controllers

import (
	services "ProjectHash/services"
	"encoding/json"
	"net/http"
)

type hashPasswordRequest struct {
	Password string `json:"password"`
}

type hashPasswordResponse struct {
	HashedPassword string `json:"hashed_password,omitempty"`
	Err            string `json:"error,omitempty"`
}

func HashPasswordHandler(svc services.PasswordService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request hashPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		hashedPassword, err := svc.HashPassword(request.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		response := hashPasswordResponse{
			HashedPassword: hashedPassword,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}