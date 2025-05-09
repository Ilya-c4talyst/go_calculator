package auth_grpc

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/gorm"

	pb "github.com/Ilya-c4talyst/go_calculator/auth/proto"
	"github.com/Ilya-c4talyst/go_calculator/auth/utils"
)

// Сервер авторизации
type authServer struct {
	pb.UnimplementedAuthServiceServer
	db *gorm.DB
}

// Функция для gRPC
func (s *authServer) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		return &pb.TokenResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &pb.TokenResponse{
		Valid:  true,
		UserId: uint32(claims.UserID),
	}, nil
}

// Старт gRPC
func StartGRPCServer(db *gorm.DB) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &authServer{db: db})

	log.Println("gRPC server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
