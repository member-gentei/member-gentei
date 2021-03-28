package clients

import (
	"context"

	"cloud.google.com/go/firestore"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// NewRetryFirestoreClient implements a client that performs 50ms +/- 10% linear backoff when it gets
// an 'Unavailable' gRPC return code.
func NewRetryFirestoreClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	unavailable := grpc_retry.WithCodes(codes.Unavailable)
	retryStreamInterceptor := grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(unavailable))
	retryUnaryInterceptor := grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(unavailable))
	return firestore.NewClient(
		ctx, projectID,
		option.WithGRPCDialOption(retryStreamInterceptor),
		option.WithGRPCDialOption(retryUnaryInterceptor),
	)
}

func NewYoutubeService(ctx context.Context) (*youtube.Service, error) {
	return nil, nil
}
