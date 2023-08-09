package services

import (
	"golang.org/x/crypto/bcrypt"
)

// interface that defines methods for password hashing and validation
type PasswordService interface {
	HashPassword(password string) (string, error)
	ValidatePassword(hashedPassword, password string) (bool, error)
}

// implements the PasswordService interface
type passwordService struct{}

// creates a new instance of the PasswordService
func NewPasswordService() PasswordService {
	return &passwordService{}
}

// hashes the provided password using bcrypt
func (s *passwordService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// validates a password against a hashed password using bcrypt
func (s *passwordService) ValidatePassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
