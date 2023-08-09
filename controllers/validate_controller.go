package controllers

import (
	services "ProjectHash/services"
	"encoding/json"
	"net/http"
)

// represents the JSON request payload structure for validating a password
type validatePasswordRequest struct {
	HashedPassword string `json:"hashed_password"`
	Password       string `json:"password"`
}

// epresents the JSON response structure for the password validation result
type validatePasswordResponse struct {
	Valid bool   `json:"valid"`
	Err   string `json:"error,omitempty"`
}

// handles the HTTP request to validate a password using the provided service
func ValidatePasswordHandler(svc services.PasswordService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//decode the incoming JSON request into validatePasswordRequest struct
		var request validatePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		//validate the password using the injected service
		valid, err := svc.ValidatePassword(request.HashedPassword, request.Password)
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		//create a response object with the validation result
		response := validatePasswordResponse{
			Valid: valid,
		}

		//set the response content type and encode the response JSON.
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
