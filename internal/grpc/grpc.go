package grpc

import (
	"fmt"
	ssov1 "github.com/aaanger/proto-1/gen/go/sso"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client struct {
	api *ssov1.AuthClient
}

func NewClient(addr string, timeout time.Duration, retriesCount int) (*Client, error) {
	retryOpts := []retry.CallOption{
		retry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		retry.WithMax(uint(retriesCount)),
		retry.WithPerRetryTimeout(timeout),
	}

	cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(retry.UnaryClientInterceptor(retryOpts...)))

	if err != nil {
		return nil, fmt.Errorf("new grpc client: %w", err)
	}

	grpcClient := ssov1.NewAuthClient(cc)

	return &Client{api: &grpcClient}, nil
}
