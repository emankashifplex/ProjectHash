package controllers

import (
	"ProjectHash/services"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHashPasswordHandler(t *testing.T) {
	// Create an instance of the actual service
	actualService := services.NewPasswordService()

	// Prepare the request
	requestBody := []byte(`{"password": "secretpassword"}`)
	req, err := http.NewRequest("POST", "/hash", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a recorder to capture the HTTP response
	rr := httptest.NewRecorder()

	// Create an HTTP handler using the actual service
	handler := http.HandlerFunc(HashPasswordHandler(actualService))

	// Serve the HTTP request and capture the response
	handler.ServeHTTP(rr, req)

	// Check if the HTTP status code is as expected
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	// Define the expected response structure
	expectedResponse := `{"hashed_password":`

	// Check if the response body contains the expected response structure
	if !bytes.Contains(rr.Body.Bytes(), []byte(expectedResponse)) {
		t.Errorf("Expected response body to contain %s, but got %s", expectedResponse, rr.Body.String())
	}
}
