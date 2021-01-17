package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/didip/tollbooth/v5"
	"github.com/didip/tollbooth/v5/limiter"
)

type FunctionHooks map[string][]string

// GothicMicroserviceClaims gothic's JWT claims
type GothicMicroserviceClaims struct {
	jwt.StandardClaims
	SiteURL       string        `json:"site_url"`
	FunctionHooks FunctionHooks `json:"function_hooks"`
}

// UnmarshalJSON unmarshal from JSON
func (f *FunctionHooks) UnmarshalJSON(b []byte) error {
	var raw map[string][]string
	err := json.Unmarshal(b, &raw)
	if err == nil {
		*f = FunctionHooks(raw)
		return nil
	}
	// If unmarshaling into map[string][]string fails, try legacy format.
	var legacy map[string]string
	err = json.Unmarshal(b, &legacy)
	if err != nil {
		return err
	}
	if *f == nil {
		*f = make(FunctionHooks)
	}
	for event, hook := range legacy {
		(*f)[event] = []string{hook}
	}
	return nil
}

func addGetBody(w http.ResponseWriter, req *http.Request) (context.Context, error) {
	if req.Method == http.MethodGet {
		return req.Context(), nil
	}

	if req.Body == nil || req.Body == http.NoBody {
		return nil, badRequestError("request must provide a body")
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, internalServerError("Error reading body").WithInternalError(err)
	}
	req.GetBody = func() (io.ReadCloser, error) {
		return ioutil.NopCloser(bytes.NewReader(buf)), nil
	}
	req.Body, _ = req.GetBody()
	return req.Context(), nil
}

func (a *API) limitHandler(lmt *limiter.Limiter) middlewareHandler {
	return func(w http.ResponseWriter, req *http.Request) (context.Context, error) {
		c := req.Context()
		if limitHeader := a.config.RateLimit; limitHeader != "" {
			key := req.Header.Get(a.config.RateLimit)
			err := tollbooth.LimitByKeys(lmt, []string{key})
			if err != nil {
				return c, httpError(http.StatusTooManyRequests, "rate limit exceeded")
			}
		}
		return c, nil
	}
}

func (a *API) requireAdminCredentials(w http.ResponseWriter, req *http.Request) (context.Context, error) {
	t, err := a.extractBearerToken(w, req)
	if err != nil {
		return nil, err
	} else if t == "" {
		return nil, errors.New("bearer token not found")
	}

	c, err := a.parseJWTClaims(t, req, w)
	if err != nil {
		return nil, err
	}

	return a.requireAdmin(c, w, req)
}

func (a *API) requireEmailProvider(w http.ResponseWriter, req *http.Request) (context.Context, error) {
	ctx := req.Context()
	config := a.config

	if config.External.Email.Disabled {
		return nil, badRequestError("Unsupported email provider")
	}

	return ctx, nil
}

func (a *API) requireConfirmed(w http.ResponseWriter, req *http.Request) (context.Context, error) {
	ctx, err := a.requireAuthentication(w, req)
	if err != nil {
		return nil, err
	}
	claims := getClaims(ctx)
	if claims == nil {
		return nil, errors.New("invalid token")
	}
	if !claims.Confirmed {
		return nil, errors.New("unconfirmed user")
	}
	return ctx, nil
}
