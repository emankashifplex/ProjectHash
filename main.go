package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	pb "ProjectHash/pb"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type PasswordService interface {
	HashPassword(password string) (string, error)
	ValidatePassword(hashedPassword, password string) error
}

type passwordService struct{}

func (s *passwordService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *passwordService) ValidatePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type hashPasswordRequest struct {
	Password string `json:"password"`
}

type hashPasswordResponse struct {
	HashedPassword string `json:"hashed_password,omitempty"`
	Err            string `json:"error,omitempty"`
}

type validatePasswordRequest struct {
	HashedPassword string `json:"hashed_password"`
	Password       string `json:"password"`
}

type validatePasswordResponse struct {
	Valid bool   `json:"valid"`
	Err   string `json:"error,omitempty"`
}

func hashPasswordHandler(svc PasswordService) http.HandlerFunc {
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

func validatePasswordHandler(svc PasswordService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request validatePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err := svc.ValidatePassword(request.HashedPassword, request.Password)
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		response := validatePasswordResponse{
			Valid: true,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

// GRPC Server implementation
type passwordServiceServer struct {
	svc PasswordService
	pb.UnimplementedPasswordServiceServer
}

func (s *passwordServiceServer) HashPassword(ctx context.Context, req *pb.HashPasswordRequest) (*pb.HashPasswordResponse, error) {
	hashedPassword, err := s.svc.HashPassword(req.Password)
	if err != nil {
		return &pb.HashPasswordResponse{
			HashedPassword: "",
			Error:          err.Error(),
		}, nil
	}
	return &pb.HashPasswordResponse{
		HashedPassword: hashedPassword,
		Error:          "",
	}, nil
}

func (s *passwordServiceServer) ValidatePassword(ctx context.Context, req *pb.ValidatePasswordRequest) (*pb.ValidatePasswordResponse, error) {
	err := s.svc.ValidatePassword(req.HashedPassword, req.Password)
	if err != nil {
		return &pb.ValidatePasswordResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}
	return &pb.ValidatePasswordResponse{
		Valid: true,
		Error: "",
	}, nil
}

func main() {

	svc := &passwordService{}

	//HTTP Server Go routine
	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/hash", hashPasswordHandler(svc)).Methods("POST")
		router.HandleFunc("/validate", validatePasswordHandler(svc)).Methods("POST")

		httpAddr := ":8080"
		fmt.Println("HTTP server is running on http://localhost" + httpAddr)
		err := http.ListenAndServe(httpAddr, router)
		if err != nil {
			fmt.Println("Error starting HTTP server:", err)
		}
	}()

	//GRPC Server Go routine
	go func() {
		grpcServer := grpc.NewServer()
		pb.RegisterPasswordServiceServer(grpcServer, &passwordServiceServer{svc: svc})

		reflection.Register(grpcServer)

		grpcAddr := ":9090"
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			fmt.Println("Failed to listen:", err)
			return
		}

		fmt.Println("gRPC server is running on http://localhost" + grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Println("Failed to serve:", err)
		}
	}()

	select {}

}
