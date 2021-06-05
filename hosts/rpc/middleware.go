package rpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/jwt"
	"github.com/jrapoport/gothic/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Authorization metadata
const (
	Authorization = "authorization"
	BearerScheme  = "bearer"
)

// Authenticator is a jwt authenticator
type Authenticator struct {
	c   config.JWT
	log log.Logger
}

// NewAuthenticator returns a new jwt authenticator
func NewAuthenticator(c config.JWT, log log.Logger) *Authenticator {
	log = log.WithName("grpc-jwt")
	return &Authenticator{c, log}
}

// UnaryServerInterceptor returns a unary interceptor for jwt authentication.
func (a *Authenticator) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpc_auth.UnaryServerInterceptor(a.authFunc)
}

// StreamServerInterceptor returns a stream interceptor for jwt authentication.
func (a *Authenticator) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return grpc_auth.StreamServerInterceptor(a.authFunc)
}

func (a *Authenticator) authFunc(ctx context.Context) (context.Context, error) {
	return Authenticate(ctx, a.c)
}

// Authenticate parses the jwt claims and authenticates a grpc request
func Authenticate(ctx context.Context, c config.JWT) (context.Context, error) {
	claims := GetClaims(ctx)
	if claims != nil {
		return ctx, nil
	}
	token, err := grpc_auth.AuthFromMD(ctx, BearerScheme)
	if err != nil {
		return ctx, err
	}
	claims, err = parseToken(c, token)
	if err != nil {
		return ctx, status.Error(codes.Unauthenticated, err.Error())
	}
	return WithClaims(ctx, claims), nil
}

func parseToken(c config.JWT, tok string) (jwt.Claims, error) {
	claims, err := jwt.ParseUserClaims(c, tok)
	if err != nil {
		return nil, err
	}
	if claims.Subject() == uuid.Nil.String() {
		err = errors.New("invalid user id")
		return nil, err
	}
	return claims, nil
}
