package grpcutil

import (
	"context"
	"crypto/tls"
	"errors"

	"github.com/iotexproject/iotex-core/action"

	"github.com/iotexproject/iotex-proto/golang/iotexapi"

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
func GetReceiptByActionHash(url string, secure bool, hash string) error {
	conn, err := ConnectToEndpoint(url, secure)
	if err != nil {
		return err
	}
	defer conn.Close()
	cli := iotexapi.NewAPIServiceClient(conn)

	request := iotexapi.GetReceiptByActionRequest{ActionHash: hash}
	response, err := cli.GetReceiptByAction(context.Background(), &request)
	if err != nil {
		return err
	}
	if response.ReceiptInfo.Receipt.Status != action.SuccessReceiptStatus {
		return errors.New("action fail:" + hash)
	}
	return nil
}
