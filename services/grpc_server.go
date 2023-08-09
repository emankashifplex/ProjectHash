package services

import (
	pb "ProjectHash/pb"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// a struct implementing the gRPC service interface.
type grpcPasswordService struct {
	pb.UnimplementedPasswordServiceServer                 //embed the generated gRPC interface
	svc                                   PasswordService //inject the PasswordService dependency
}

// creates a new gRPC server with the grpcPasswordService
func NewGRPCServer(svc PasswordService) *grpc.Server {
	grpcServer := grpc.NewServer()
	grpcPasswordService := &grpcPasswordService{svc: svc}
	pb.RegisterPasswordServiceServer(grpcServer, grpcPasswordService)
	return grpcServer
}

// implements the gRPC HashPassword method
func (s *grpcPasswordService) HashPassword(ctx context.Context, req *pb.HashPasswordRequest) (*pb.HashPasswordResponse, error) {
	password := req.GetPassword()

	// hash the password using the injected service
	hashedPassword, err := s.svc.HashPassword(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error hashing password: %v", err)
	}

	// create and return the response
	response := &pb.HashPasswordResponse{
		HashedPassword: hashedPassword,
	}
	return response, nil
}

// implements the gRPC ValidatePassword method
func (s *grpcPasswordService) ValidatePassword(ctx context.Context, req *pb.ValidatePasswordRequest) (*pb.ValidatePasswordResponse, error) {
	hashedPassword := req.GetHashedPassword()
	password := req.GetPassword()

	// validate the password using the injected service
	valid, err := s.svc.ValidatePassword(hashedPassword, password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid password: %v", err)
	}

	// create and return the response
	response := &pb.ValidatePasswordResponse{
		Valid: valid,
	}
	return response, nil
}
