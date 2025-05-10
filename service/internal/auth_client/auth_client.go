package auth_client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/Ilya-c4talyst/go_calculator/service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// Клиент аутентификации
type Client struct {
	conn   *grpc.ClientConn
	client pb.AuthServiceClient
}

// Создание клиента
func New(addr string) (*Client, error) {
	if addr == "" {
		return nil, fmt.Errorf("empty gRPC server address")
	}

	log.Printf("Attempting to connect to gRPC server at %s", addr)

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5 * time.Second),
		grpc.WithReturnConnectionError(),
	}

	if os.Getenv("ENV") == "dev" {
		addr = "localhost:50051"
	}

	conn, err := grpc.Dial(addr, dialOptions...)
	if err != nil {
		log.Printf("gRPC connection error details: %v", err)

		// Проверяем конкретные типы ошибок
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("gRPC connection timeout to %s", addr)
		}

		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %v", addr, err)
	}

	// Проверяем состояние соединения
	state := conn.GetState()
	log.Printf("gRPC connection state: %s", state.String())

	if state != connectivity.Ready {
		conn.Close()
		return nil, fmt.Errorf("gRPC connection not ready, state: %s", state)
	}

	log.Printf("Successfully connected to gRPC server at %s", addr)
	return &Client{
		conn:   conn,
		client: pb.NewAuthServiceClient(conn),
	}, nil
}

// Валидация токена
func (c *Client) ValidateToken(token string) (uint32, error) {
	resp, err := c.client.ValidateToken(context.Background(), &pb.TokenRequest{
		Token: token,
	})
	if err != nil {
		return 0, fmt.Errorf("auth service error: %v", err)
	}

	if !resp.Valid {
		return 0, fmt.Errorf(resp.Error)
	}

	return resp.UserId, nil
}

// Закрытие клиента
func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		log.Printf("Failed to close gRPC connection: %v", err)
	}
}
