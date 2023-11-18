package client

import (
	"github.com/romankravchuk/eldorado/internal/services/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(url string) (proto.AuthServiceClient, error) {
	con, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return proto.NewAuthServiceClient(con), nil
}
