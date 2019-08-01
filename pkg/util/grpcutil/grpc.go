package grpcutil

import (
	"crypto/tls"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func ConnectToEndpoint(url string, secure bool) (*grpc.ClientConn, error) {
	endpoint := url
	if endpoint == "" {
		return nil, errors.New(`endpoint is empty`)
	}
	if !secure {
		return grpc.Dial(endpoint, grpc.WithInsecure())
	}
	return grpc.Dial(endpoint, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
}
