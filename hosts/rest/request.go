package rest

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
	"github.com/segmentio/encoding/json"
	"github.com/sirupsen/logrus"
)

// HTTP Headers
const (
	Authorization   = "Authorization"
	BearerScheme    = "Bearer"
	ContentType     = "Content-Type"
	UseCookieHeader = "X-Use-Cookie"
	ForwardedProto  = "X-Forwarded-Proto"
)

const (
	// JSONContent json content type.
	JSONContent = "application/json"
	// SessionCookie token.
	SessionCookie = "Session"
	// JWTCookieKey jwt token key for a token.
	JWTCookieKey = key.JWT
	// JWTQueryKey jwt token key for a query string.
	JWTQueryKey = key.JWT
)

// Request for a request
type Request struct {
	Code      string        `json:"code" form:"code"`
	Provider  provider.Name `json:"provider" form:"provider"`
	ReCaptcha string        `json:"recaptcha" form:"recaptcha"`
	Sort      store.Sort    `json:"sort"  form:"sort"`
}

// FromRequest adds the authenticated request context to a context
func FromRequest(r *http.Request) context.Context {
	req := &Request{}
	err := UnmarshalRequest(r, req)
	ctx := context.WithContext(r.Context())
	ctx.SetIPAddress(r.RemoteAddr)
	ctx.SetCode(req.Code)
	ctx.SetProvider(req.Provider)
	ctx.SetReCaptcha(req.ReCaptcha)
	ctx.SetSort(req.Sort)
	c, err := GetUserClaims(r)
	if err != nil {
		return ctx
	}
	ctx.SetProvider(c.Provider)
	if c.Admin {
		ctx.SetAdminID(c.UserID())
	} else {
		ctx.SetUserID(c.UserID())
	}
	return ctx
}

type loggerKey struct{}

// WithLogger adds a logger to the context of an http request.
func WithLogger(r *http.Request, l logrus.FieldLogger) *http.Request {
	ctx := context.WithValue(r.Context(), loggerKey{}, l)
	return r.WithContext(ctx)
}

// GetLogger gets a logger to the context of an http request.
func GetLogger(r *http.Request) logrus.FieldLogger {
	log, ok := r.Context().Value(loggerKey{}).(logrus.FieldLogger)
	if !ok {
		l := logrus.New()
		l.SetLevel(logrus.PanicLevel)
		return l
	}
	return log
}

type claimsKey struct{}

// WithClaims adds a set of jwt claims to the context of an http request.
func WithClaims(r *http.Request, c jwt.Claims) *http.Request {
	ctx := context.WithValue(r.Context(), claimsKey{}, c)
	return r.WithContext(ctx)
}

// GetClaims gets a set of jwt claims from the context of an http request.
func GetClaims(r *http.Request) jwt.Claims {
	c, _ := r.Context().Value(claimsKey{}).(jwt.Claims)
	return c
}

// GetUserClaims gets the user claims from the context of an http request.
func GetUserClaims(r *http.Request) (*jwt.UserClaims, error) {
	claims := GetClaims(r)
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

// GetUserID gets the user id from the jwt claims of an http request.
func GetUserID(r *http.Request) (uuid.UUID, error) {
	claims, err := GetUserClaims(r)
	if err != nil {
		return uuid.Nil, err
	}
	uid := claims.UserID()
	if uid == uuid.Nil {
		err = errors.New("invalid user id")
		return uuid.Nil, err
	}
	return uid, nil
}

// UnmarshalRequest unmarshal a Request from an http.Request
func UnmarshalRequest(r *http.Request, v interface{}) error {
	if r.Form == nil {
		const defaultMaxMemory = 32 << 20 // 32 MB
		// it looks like it will return "no body" but still do the right
		// thing with the form. net/http/request.go ignores it also.
		_ = r.ParseMultipartForm(defaultMaxMemory)
	}
	if isJSONRequest(r) {
		var b bytes.Buffer
		_, err := io.Copy(&b, r.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(b.Bytes(), v)
		if err != nil {
			return err
		}
		// read & reset the body so it may be read again.
		r.Body = ioutil.NopCloser(&b)
		return nil
	}
	switch t := v.(type) {
	case *url.Values:
		*t = r.Form
		return nil
	case *types.Map:
		*t = utils.URLValuesToMap(r.Form, true)
		return nil
	case *store.Filters:
		*t = utils.URLValuesToMap(r.Form, true)
		return nil
	default:
		break
	}
	dec := form.NewDecoder()
	dec.RegisterCustomTypeFunc(func(val []string) (interface{}, error) {
		data := types.Map{}
		if len(val) <= 0 {
			return data, nil
		}
		err := data.Scan(val[0])
		if err != nil {
			return data, err
		}
		return data, nil
	}, types.Map{})
	return dec.Decode(v, r.Form)
}

func isJSONRequest(r *http.Request) bool {
	ct := r.Header.Get(ContentType)
	return ct == JSONContent
}

// RequestToken returns a jwt token from an http.Request
func RequestToken(r *http.Request) string {
	var t string
	find := []func(r *http.Request) string{
		headerToken,
		cookieToken,
		queryToken,
	}
	for _, fn := range find {
		if t = fn(r); t != "" {
			break
		}
	}
	return t
}

// ParseClaims parses a set of jwt claims from an http request.
func ParseClaims(r *http.Request, jc config.JWT, tok string) (*http.Request, error) {
	claims, err := jwt.ParseUserClaims(jc, tok)
	if err != nil {
		return nil, err
	}
	if claims.UserID() == uuid.Nil {
		err = errors.New("invalid user id")
		return nil, err
	}
	return WithClaims(r, claims), nil
}

func headerToken(r *http.Request) string {
	const Bearer = BearerScheme + " "
	n := len(Bearer)
	auth := r.Header.Get(Authorization)
	var bearer string
	if len(auth) >= n && strings.EqualFold(auth[:n], Bearer) {
		bearer = auth[n:]
	}
	return bearer
}

func cookieToken(r *http.Request) string {
	// token from a cookie named "jwt".
	cookie, err := r.Cookie(JWTCookieKey)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func queryToken(r *http.Request) string {
	// token from query string ("jwt").
	return r.URL.Query().Get(JWTQueryKey)
}
