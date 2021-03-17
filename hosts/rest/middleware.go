package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httptracer"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/opentracing/opentracing-go"
)

// Logger is a middleware that logs the start and end of each request.
func Logger(h http.Handler) http.Handler {
	//return middleware.Logger(h)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == config.HealthEndpoint {
			h.ServeHTTP(w, r)
			return
		}
		middleware.Logger(h).ServeHTTP(w, r)
	})
}

// CORS creates a new Cors handler.
func CORS(h http.Handler) http.Handler {
	fn := cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			Authorization,
			ContentType,
			"X-CSRF-AccessToken",
			"X-JWT-AUD",
			UseCookieHeader,
		},
		AllowCredentials: true,
		MaxAge:           300,
	})
	return fn(h)
}

// Tracer creates a new Tracer handler.
func Tracer(h http.Handler) http.Handler {
	tr := opentracing.GlobalTracer()
	if tr == nil {
		return h
	}
	fn := httptracer.Tracer(tr, httptracer.Config{
		SkipFunc: func(r *http.Request) bool {
			return r.URL.Path == config.HealthEndpoint
		},
		Tags: map[string]interface{}{
			// datadog, turn on metrics for http.request stats
			"_dd.measured": 1,
			// datadog, event sample rate
			"_dd1.sr.eausr": 1,
		},
	})
	return fn(h)
}

// Authenticator creates a new JWT handler with passed options.
func Authenticator(jc config.JWT) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// did we already do this?
			c := GetClaims(r)
			if c != nil {
				// yes, skip everything else
				next.ServeHTTP(w, r)
				return
			}
			l := GetLogger(r)
			tok := RequestToken(r)
			if tok == "" {
				l.Error("token not found")
				ResponseCode(w, http.StatusUnauthorized, nil)
				return
			}
			r, err := ParseClaims(r, jc, tok)
			if err != nil {
				l.Error(err)
				ResponseCode(w, http.StatusUnauthorized, nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// AdminUser creates a new JWT handler that checks for admin permissions.
func AdminUser(next http.Handler) http.Handler {
	return claimCheck(func(claims jwt.UserClaims) error {
		if !claims.Confirmed {
			return errors.New("user not confirmed")
		}
		if !claims.Admin {
			return errors.New("admin access required")
		}
		return nil
	})(next)
}

// ConfirmedUser creates a new JWT handler that checks for user confirmation.
func ConfirmedUser(next http.Handler) http.Handler {
	return claimCheck(func(claims jwt.UserClaims) error {
		if !claims.Confirmed {
			return errors.New("user not confirmed")
		}
		return nil
	})(next)
}

func claimCheck(check func(claims jwt.UserClaims) error) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := GetLogger(r)
			c := GetClaims(r)
			if c == nil {
				l.Error("jwt claims not found")
				ResponseCode(w, http.StatusUnauthorized, nil)
				return
			}
			uc, ok := c.(jwt.UserClaims)
			if !ok {
				l.Error("user claims not found")
				return
			}
			if err := check(uc); err != nil {
				l.Error(err)
				ResponseCode(w, http.StatusForbidden, nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// UseCookie adds a jwt token to a cookie.
func UseCookie(w http.ResponseWriter, r *http.Request, token string, exp time.Duration) {
	use := r.Header.Get(UseCookieHeader)
	if use == "" {
		return
	}
	cookie := &http.Cookie{
		Name:     JWTCookieKey,
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	if use != SessionCookie {
		cookie.Expires = time.Now().UTC().Add(exp)
		cookie.MaxAge = int(exp / time.Second)
	}
	http.SetCookie(w, cookie)
}

// ClearCookie clears a jwt token to a cookie.
func ClearCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     JWTCookieKey,
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}
	http.SetCookie(w, c)
}
