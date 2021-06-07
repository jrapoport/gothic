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

	IPAddress() string
	SetIPAddress(string)

	Provider() provider.Name
	SetProvider(provider.Name)

	UserID() uuid.UUID
	SetUserID(uuid.UUID)

	AdminID() uuid.UUID
	SetAdminID(uuid.UUID)

	IsAdmin() bool

	Code() string
	SetCode(string)

	ReCaptcha() string
	SetReCaptcha(string)

	Sort() store.Sort
	SetSort(store.Sort)
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

type ipKey struct{}

func (ctx apiContext) IPAddress() string {
	v, _ := ctx.Value(ipKey{}).(string)
	return v
}

func (ctx *apiContext) SetIPAddress(ip string) {
	if ip == "" {
		return
	}
	ctx.setValue(ipKey{}, ip)
}

type providerKey struct{}

func (ctx apiContext) Provider() provider.Name {
	v, _ := ctx.Value(providerKey{}).(provider.Name)
	return v
}

func (ctx *apiContext) SetProvider(s provider.Name) {
	if s == provider.Unknown {
		return
	}
	ctx.setValue(providerKey{}, s)
}

type uidKey struct{}

func (ctx apiContext) UserID() uuid.UUID {
	v, _ := ctx.Value(uidKey{}).(uuid.UUID)
	return v
}

func (ctx *apiContext) SetUserID(uid uuid.UUID) {
	if uid == uuid.Nil {
		return
	}
	ctx.setValue(uidKey{}, uid)
}

type aidKey struct{}

func (ctx apiContext) AdminID() uuid.UUID {
	v, _ := ctx.Value(aidKey{}).(uuid.UUID)
	return v
}

func (ctx *apiContext) SetAdminID(uid uuid.UUID) {
	if uid == uuid.Nil {
		return
	}
	ctx.setValue(aidKey{}, uid)
}

func (ctx apiContext) IsAdmin() bool {
	return ctx.AdminID() != uuid.Nil
}

type codeKey struct{}

func (ctx apiContext) Code() string {
	v, _ := ctx.Value(codeKey{}).(string)
	return v
}

func (ctx *apiContext) SetCode(code string) {
	if code == "" {
		return
	}
	ctx.setValue(codeKey{}, code)
}

type recaptchaKey struct{}

func (ctx apiContext) ReCaptcha() string {
	v, _ := ctx.Value(recaptchaKey{}).(string)
	return v
}

func (ctx *apiContext) SetReCaptcha(tok string) {
	if tok == "" {
		return
	}
	ctx.setValue(recaptchaKey{}, tok)
}

type sortKey struct{}

func (ctx apiContext) Sort() store.Sort {
	v, _ := ctx.Value(sortKey{}).(store.Sort)
	return v
}

func (ctx *apiContext) SetSort(s store.Sort) {
	ctx.setValue(sortKey{}, s)
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
