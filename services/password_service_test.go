package services

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	svc := NewPasswordService()

	password := "secretpassword"
	hashedPassword, err := svc.HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Error("Hashed password is empty")
	}
}

func TestValidatePassword(t *testing.T) {
	svc := NewPasswordService()

	password := "secretpassword"
	hashedPassword, _ := svc.HashPassword(password)

	valid, err := svc.ValidatePassword(hashedPassword, password)
	if err != nil {
		t.Errorf("Error validating password: %v", err)
	}

	if !valid {
		t.Error("Password validation failed")
	}
}
