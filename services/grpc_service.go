package services

import (
	"ProjectHash/pb"
	"context"
)

type GRPCService struct {
	pb.UnimplementedPasswordServiceServer
	svc PasswordService
}

func NewGRPCService(svc PasswordService) pb.PasswordServiceServer {
	return &GRPCService{svc: svc}
}

func (s *GRPCService) HashPassword(ctx context.Context, req *pb.HashPasswordRequest) (*pb.HashPasswordResponse, error) {
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

func (s *GRPCService) ValidatePassword(ctx context.Context, req *pb.ValidatePasswordRequest) (*pb.ValidatePasswordResponse, error) {
	valid, err := s.svc.ValidatePassword(req.HashedPassword, req.Password)
	if err != nil {
		return &pb.ValidatePasswordResponse{
			Valid: valid,
			Error: err.Error(),
		}, nil
	}
	return &pb.ValidatePasswordResponse{
		Valid: valid,
		Error: "",
	}, nil
}
