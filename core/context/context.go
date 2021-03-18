package context

import (
	"context"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store"
)

// Context is for api calls.
type Context interface {
	context.Context

	GetCode() string
	GetIPAddress() string
	GetProvider() provider.Name
	GetReCaptcha() string
	GetSort() store.Sort
	GetUserID() uuid.UUID
	GetAdminID() uuid.UUID

	SetCode(string)
	SetIPAddress(string)
	SetProvider(provider.Name)
	SetReCaptcha(string)
	SetSort(store.Sort)
	SetUserID(uuid.UUID)
	SetAdminID(uuid.UUID)
}

// Background returns a new wrapped background context.
func Background() Context {
	return &apiContext{context.Background()}
}

// WithValue wraps context.WithValue
func WithValue(parent context.Context, key, val interface{}) Context {
	return &apiContext{context.WithValue(parent, key, val)}
}

type apiContext struct {
	context.Context
}

var _ Context = (*apiContext)(nil)

type codeKey struct{}
type ipKey struct{}
type providerKey struct{}
type recaptchaKey struct{}
type sortKey struct{}
type uidKey struct{}
type aidKey struct{}

func (ctx apiContext) GetCode() string {
	tok, _ := ctx.Value(codeKey{}).(string)
	return tok
}

func (ctx apiContext) GetIPAddress() string {
	e, _ := ctx.Value(ipKey{}).(string)
	return e
}

func (ctx apiContext) GetProvider() provider.Name {
	e, _ := ctx.Value(providerKey{}).(provider.Name)
	return e
}

func (ctx apiContext) GetReCaptcha() string {
	tok, _ := ctx.Value(recaptchaKey{}).(string)
	return tok
}

func (ctx apiContext) GetSort() store.Sort {
	s, _ := ctx.Value(sortKey{}).(store.Sort)
	return s
}

func (ctx apiContext) GetUserID() uuid.UUID {
	e, _ := ctx.Value(uidKey{}).(uuid.UUID)
	return e
}

func (ctx apiContext) GetAdminID() uuid.UUID {
	e, _ := ctx.Value(aidKey{}).(uuid.UUID)
	return e
}

func (ctx *apiContext) SetCode(code string) {
	if code == "" {
		return
	}
	ctx.setValue(codeKey{}, code)
}

func (ctx *apiContext) SetIPAddress(ip string) {
	if ip == "" {
		return
	}
	ctx.setValue(ipKey{}, ip)
}

func (ctx *apiContext) SetProvider(s provider.Name) {
	if s == provider.Unknown {
		return
	}
	ctx.setValue(providerKey{}, s)
}

func (ctx *apiContext) SetReCaptcha(tok string) {
	if tok == "" {
		return
	}
	ctx.setValue(recaptchaKey{}, tok)
}

func (ctx *apiContext) SetSort(s store.Sort) {
	if s == "" {
		return
	}
	ctx.setValue(sortKey{}, s)
}

func (ctx *apiContext) SetUserID(uid uuid.UUID) {
	if uid == uuid.Nil {
		return
	}
	ctx.setValue(uidKey{}, uid)
}

func (ctx *apiContext) SetAdminID(uid uuid.UUID) {
	if uid == uuid.Nil {
		return
	}
	ctx.setValue(aidKey{}, uid)
}

func (ctx *apiContext) setValue(key, val interface{}) {
	if ctx.Value(key) == val {
		return
	}
	*ctx = apiContext{context.WithValue(ctx.Context, key, val)}
}

// WithContext returns the request context with a context
func WithContext(ctx context.Context) Context {
	if ctx == nil {
		ctx = context.Background()
	}
	v, ok := ctx.(Context)
	if ok {
		return v
	}
	return &apiContext{ctx}
}
