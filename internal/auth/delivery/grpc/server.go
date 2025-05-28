package grpc

import (
	"context"

	"github.com/chizheg/forum/internal/auth/domain"
	pb "github.com/chizheg/forum/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	service domain.Service
	logger  *zap.Logger
}

func NewAuthServer(service domain.Service, logger *zap.Logger) *AuthServer {
	return &AuthServer{
		service: service,
		logger:  logger,
	}
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	token, err := s.service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		s.logger.Error("failed to register user", zap.Error(err))
		return &pb.RegisterResponse{
			Success: false,
			Error:   err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		Success: true,
		Token:   token,
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.service.Login(req.Username, req.Password)
	if err != nil {
		s.logger.Error("failed to login user", zap.Error(err))
		return &pb.LoginResponse{
			Success: false,
			Error:   err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.LoginResponse{
		Success: true,
		Token:   token,
	}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	userID, err := s.service.ValidateToken(req.Token)
	if err != nil {
		s.logger.Error("failed to validate token", zap.Error(err))
		return &pb.ValidateTokenResponse{
			Valid: false,
			Error: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: int32(userID),
	}, nil
}
