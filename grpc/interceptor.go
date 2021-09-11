package grpc

import (
	"context"

	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func authorizationInterceptor(authToken string) grpcauth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		token, err := grpcauth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		if token != authToken {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %q", token)
		}

		return ctx, nil
	}
}
