package rpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// Authorization rpc metadata.
	Authorization = "authorization"
	// BearerScheme rpc metadata.
	BearerScheme = "bearer"
)

// Authenticator is a jwt authenticator
type Authenticator struct {
	c   config.JWT
	log logrus.FieldLogger
}

// NewAuthenticator returns a new jwt authenticator
func NewAuthenticator(c config.JWT, log logrus.FieldLogger) *Authenticator {
	log = log.WithField("grpc-middleware", "jwt")
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

func (a *Authenticator) parseToken(tok string) (jwt.Claims, error) {
	claims, err := jwt.ParseUserClaims(a.c, tok)
	if err != nil {
		return nil, err
	}
	if claims.Subject == uuid.Nil.String() {
		err = errors.New("invalid user id")
		return nil, err
	}
	return claims, nil
}

func (a *Authenticator) authFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, BearerScheme)
	if err != nil {
		return nil, err
	}
	claims, err := a.parseToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return WithClaims(ctx, claims), nil
}
