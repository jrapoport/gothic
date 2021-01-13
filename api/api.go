package api

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/didip/tollbooth/v5"
	"github.com/didip/tollbooth/v5/limiter"
	"github.com/go-chi/chi"
	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/mailer"
	"github.com/jrapoport/gothic/storage"
	"github.com/jrapoport/gothic/util"
	"github.com/rs/cors"
	"github.com/sebest/xff"
	"github.com/sirupsen/logrus"
)

const (
	audHeaderName = "X-JWT-AUD"
)

var bearerRegexp = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

// API is the main REST API
type API struct {
	handler http.Handler
	db      *storage.Connection
	config  *conf.Configuration
}

// ListenAndServeREST starts the REST API
// let's wrap this instead
func ListenAndServeREST(a *API, globalConfig *conf.Configuration) {
	go func() {
		addr := fmt.Sprintf("%v:%v", globalConfig.Host, globalConfig.RestPort)
		logrus.Infof("Gothic REST API started on: %s", addr)
		a.ListenAndServe(addr)
	}()
}

func (a *API) ListenAndServe(hostAndPort string) {
	log := logrus.WithField("component", "api")
	server := &http.Server{
		Addr:    hostAndPort,
		Handler: a.handler,
	}

	done := make(chan struct{})
	defer close(done)
	go func() {
		util.WaitForTermination(log, done)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		server.Shutdown(ctx)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.WithError(err).Fatal("http server listen failed")
	}
}

// NewAPI creates a new REST API
func NewAPI(globalConfig *conf.Configuration, db *storage.Connection) *API {
	api := &API{config: globalConfig, db: db}

	xffmw, _ := xff.Default()
	logger := newStructuredLogger(logrus.StandardLogger())

	r := newRouter()
	r.UseBypass(xffmw.Handler)
	r.Use(addRequestID(globalConfig))
	r.Use(recoverer)
	r.UseBypass(tracer)

	r.Get("/health", api.handleHealthCheck)

	r.Route("/callback", func(r *router) {
		r.UseBypass(logger)
		r.Use(api.loadOAuthState)
		r.Get("/", api.ExternalProviderCallback)
	})

	r.Route("/", func(r *router) {
		r.UseBypass(logger)

		r.Get("/settings", api.handleSettings)

		r.Get("/authorize", api.ExternalProviderRedirect)

		r.With(api.requireAdminCredentials).Post("/invite", api.Invite)

		r.With(api.requireEmailProvider).Post("/signup", api.Signup)
		r.With(api.requireEmailProvider).Post("/recover", api.Recover)
		r.With(api.requireEmailProvider).With(api.limitHandler(
			// Allow requests at a rate of 30 per 5 minutes.
			tollbooth.NewLimiter(30.0/(60*5), &limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Hour,
			}).SetBurst(30),
		)).Post("/token", api.Token)
		r.Post("/verify", api.Verify)

		r.With(api.requireAuthentication).Post("/logout", api.Logout)

		r.Route("/user", func(r *router) {
			r.Use(api.requireAuthentication)
			r.Get("/", api.UserGet)
			r.Put("/", api.UserUpdate)
		})

		r.Route("/admin", func(r *router) {
			r.Use(api.requireAdminCredentials)

			r.Route("/audit", func(r *router) {
				r.Get("/", api.adminAuditLog)
			})

			r.Route("/users", func(r *router) {
				r.Get("/", api.adminUsers)
				r.With(api.requireEmailProvider).Post("/", api.adminUserCreate)

				r.Route("/{user_id}", func(r *router) {
					r.Use(api.loadUser)

					r.Get("/", api.adminUserGet)
					r.Put("/", api.adminUserUpdate)
					r.Delete("/", api.adminUserDelete)
				})
			})
		})

		r.Route("/saml", func(r *router) {
			r.Route("/acs", func(r *router) {
				r.Use(api.loadSAMLState)
				r.Post("/", api.ExternalProviderCallback)
			})

			r.Get("/metadata", api.SAMLMetadata)
		})
	})

	corsHandler := cors.New(cors.Options{
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", audHeaderName, useCookieHeader},
		AllowCredentials: true,
	})

	ctx := withConfig(context.Background(), globalConfig)

	api.handler = corsHandler.Handler(chi.ServerBaseContext(ctx, r))
	return api
}

func WithConfig(ctx context.Context, config *conf.Configuration) (context.Context, error) {
	ctx = withConfig(ctx, config)
	return ctx, nil
}

func (a *API) Mailer(ctx context.Context) mailer.Mailer {
	config := a.getConfig(ctx)
	return mailer.NewMailer(config)
}

func (a *API) getConfig(context.Context) *conf.Configuration {
	return a.config
}
