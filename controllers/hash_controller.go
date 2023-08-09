package controllers

import (
	services "ProjectHash/services"
	"encoding/json"
	"net/http"
)

// represents JSON request payload structure for hashing password
type hashPasswordRequest struct {
	Password string `json:"password"`
}

// represents JSON response structure for the hashed password
type hashPasswordResponse struct {
	HashedPassword string `json:"hashed_password,omitempty"`
	Err            string `json:"error,omitempty"`
}

// handles the HTTP request to hash password using the provided service
func HashPasswordHandler(svc services.PasswordService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//decode the incoming JSON request into hashPasswordRequest struct.
		var request hashPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		//hash the password using the injected service.
		hashedPassword, err := svc.HashPassword(request.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		//Create a response object with the hashed password.
		response := hashPasswordResponse{
			HashedPassword: hashedPassword,
		}

		//Set the response content type and encode the response JSON.
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
