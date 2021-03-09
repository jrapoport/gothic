package rpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	core_ctx "github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	// RealIP is the context metadata key for real ip
	RealIP = "X-Real-IP"
	// ForwardedFor is the context metadata key for forwarded for ip
	ForwardedFor = "X-Forwarded-For"
	// ReCaptchaToken is the context metadata key for a recaptcha token
	ReCaptchaToken = "X-ReCaptcha-Token"
)

// RequestContext adds the request context to a context
func RequestContext(ctx context.Context) core_ctx.Context {
	rtx := core_ctx.WithContext(ctx)
	rtx.SetIPAddress(GetRemoteIP(rtx))
	if c, err := GetUserClaims(rtx); err == nil {
		rtx.SetProvider(c.Provider)
		if c.Admin {
			rtx.SetAdminID(c.UserID())
		}
		rtx.SetUserID(c.UserID())
	}
	md, ok := metadata.FromIncomingContext(rtx)
	if !ok {
		return rtx
	}
	if v := getMetadata(md, ReCaptchaToken); v != "" {
		rtx.SetReCaptcha(v)
	}
	return rtx
}

// GetRemoteIP returns the remote ip from the context.
func GetRemoteIP(ctx context.Context) string {
	ip := realIP(ctx)
	if ip != "" {
		return ip
	}
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	return pr.Addr.String()
}

func realIP(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if xrip := getMetadata(md, RealIP); xrip != "" {
		return xrip
	}
	if xff := getMetadata(md, ForwardedFor); xff != "" {
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		return xff[:i]
	}
	return ""
}

func getMetadata(md metadata.MD, k string) string {
	for _, v := range md.Get(k) {
		if v == "" {
			continue
		}
		return v
	}
	return ""
}

type claimsKey struct{}

// WithClaims adds a set of jwt claims to the grpc context.
func WithClaims(ctx context.Context, c jwt.Claims) context.Context {
	return context.WithValue(ctx, claimsKey{}, c)
}

// GetClaims gets a set of jwt claims from the grpc context.
func GetClaims(ctx context.Context) jwt.Claims {
	c, _ := ctx.Value(claimsKey{}).(jwt.Claims)
	return c
}

// GetUserID gets the user id from the jwt claims of the grpc context.
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	uc, err := GetUserClaims(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	uid, err := uuid.Parse(uc.Subject)
	if err != nil {
		err = fmt.Errorf("invalid user id %s: %w",
			uc.Subject, err)
		return uuid.Nil, err
	}
	if uid == uuid.Nil {
		err = fmt.Errorf("invalid user id: %s", uid)
		return uuid.Nil, err
	}
	return uid, nil
}

// GetUserClaims gets the user claims from the jwt claims of the grpc context.
func GetUserClaims(ctx context.Context) (*jwt.UserClaims, error) {
	claims := GetClaims(ctx)
	if claims == nil {
		err := errors.New("jwt claims not found")
		return nil, err
	}
	uc, ok := claims.(jwt.UserClaims)
	if !ok {
		err := errors.New("invalid jwt user claims")
		return nil, err
	}
	return &uc, nil
}
