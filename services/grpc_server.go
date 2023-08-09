package services

import (
	pb "ProjectHash/pb"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcPasswordService struct {
	pb.UnimplementedPasswordServiceServer
	svc PasswordService
}

func NewGRPCServer(svc PasswordService) *grpc.Server {
	grpcServer := grpc.NewServer()
	grpcPasswordService := &grpcPasswordService{svc: svc}
	pb.RegisterPasswordServiceServer(grpcServer, grpcPasswordService)
	return grpcServer
}

func (s *grpcPasswordService) HashPassword(ctx context.Context, req *pb.HashPasswordRequest) (*pb.HashPasswordResponse, error) {
	password := req.GetPassword()

	hashedPassword, err := s.svc.HashPassword(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error hashing password: %v", err)
	}

	response := &pb.HashPasswordResponse{
		HashedPassword: hashedPassword,
	}
	return response, nil
}

func (s *grpcPasswordService) ValidatePassword(ctx context.Context, req *pb.ValidatePasswordRequest) (*pb.ValidatePasswordResponse, error) {
	hashedPassword := req.GetHashedPassword()
	password := req.GetPassword()

	valid, err := s.svc.ValidatePassword(hashedPassword, password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid password: %v", err)
	}

	response := &pb.ValidatePasswordResponse{
		Valid: valid,
	}
	return response, nil
}
