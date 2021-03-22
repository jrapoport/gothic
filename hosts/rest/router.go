package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jrapoport/gothic/config"
	"github.com/sirupsen/logrus"
)

// Root is the base route
const Root = "/"

// Router an http router
type Router struct {
	chi    chi.Router
	config *config.Config
}

// NewRouter returns a new configured router.
func NewRouter(c *config.Config) *Router {
	r := &Router{chi: chi.NewRouter(), config: c}
	r.UseDefaults()
	return r
}

// UseDefaults applies the default router middlewares.
func (r *Router) UseDefaults() {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(CORS)
	r.Use(Tracer)
	r.Use(Logger)
	r.Use(middleware.Recoverer)
	r.chi.NotFound(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, r.config.SiteURL, http.StatusSeeOther)
	})
}

// Mount attaches another http.Handler along ./pattern/*
func (r *Router) Mount(pattern string, h http.Handler) {
	r.chi.Mount(pattern, h)
}

// Route mounts a sub-Router along a `pattern`` string.
func (r *Router) Route(pattern string, fn func(*Router)) {
	r.chi.Route(pattern, func(cr chi.Router) {
		fn(&Router{chi: cr, config: r.config})
	})
}

// Handler adds routes for `pattern` that matches all HTTP methods.
func (r *Router) Handler(pattern string, fn http.HandlerFunc) {
	r.chi.HandleFunc(pattern, fn)
}

// Get HTTP-method routing along `pattern`
func (r *Router) Get(pattern string, fn http.HandlerFunc) {
	r.chi.Get(pattern, fn)
}

// Post HTTP-method routing along `pattern`
func (r *Router) Post(pattern string, fn http.HandlerFunc) {
	r.chi.Post(pattern, fn)
}

// Put HTTP-method routing along `pattern`
func (r *Router) Put(pattern string, fn http.HandlerFunc) {
	r.chi.Put(pattern, fn)
}

// Delete HTTP-method routing along `pattern`
func (r *Router) Delete(pattern string, fn http.HandlerFunc) {
	r.chi.Delete(pattern, fn)
}

// With adds inline middlewares for an endpoint handler.
func (r *Router) With(middlewares ...func(http.Handler) http.Handler) *Router {
	cr := r.chi.With(middlewares...)
	return &Router{chi: cr, config: r.config}
}

// Use appends one or more middlewares onto the Router stack.
func (r *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	r.chi.Use(middlewares...)
}

// ServeHTTP should write reply headers and data to the ResponseWriter.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chi.ServeHTTP(w, req)
}

// UseLogger sets the logger to use for the Router stack.
func (r *Router) UseLogger(log logrus.FieldLogger) {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = WithLogger(r, log)
			next.ServeHTTP(w, r)
		})
	})
}

// RateLimit adds inline middlewares to rate limit endpoint handlers.
func (r *Router) RateLimit() *Router {
	if r.config == nil || r.config.RateLimit <= 0 {
		return r
	}
	return r.With(httprate.LimitByIP(100, r.config.RateLimit))
}

// Authenticated adds inline middlewares to enforce jwt for endpoint handlers.
func (r *Router) Authenticated() *Router {
	if r.config == nil {
		return r
	}
	return r.With(Authenticator(r.config.JWT))
}

// Authenticate adds inline middlewares to enforce
// jwt for endpoint handlers with the config.
func (r *Router) Authenticate(c config.JWT) *Router {
	return r.With(Authenticator(c))
}

// Admin adds inline middlewares to enforce
// admin permissions for endpoint handlers.
func (r *Router) Admin() *Router {
	return r.With(AdminUser)
}

// Confirmed adds inline middlewares to enforce
// confirmed accounts permissions for endpoint handlers
func (r *Router) Confirmed() *Router {
	return r.With(ConfirmedUser)
}

// URLParam returns the url parameter from a http.Request object.
func URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
