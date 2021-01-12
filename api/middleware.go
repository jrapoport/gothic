package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/didip/tollbooth/v5"
	"github.com/didip/tollbooth/v5/limiter"
)

type FunctionHooks map[string][]string

type GothicMicroserviceClaims struct {
	jwt.StandardClaims
	SiteURL       string        `json:"site_url"`
	InstanceID    string        `json:"id"`
	FunctionHooks FunctionHooks `json:"function_hooks"`
}

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
		if limitHeader := a.config.RateLimitHeader; limitHeader != "" {
			key := req.Header.Get(a.config.RateLimitHeader)
			err := tollbooth.LimitByKeys(lmt, []string{key})
			if err != nil {
				return c, httpError(http.StatusTooManyRequests, "Rate limit exceeded")
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
	config := a.getConfig(ctx)

	if config.External.Email.Disabled {
		return nil, badRequestError("Unsupported email provider")
	}

	return ctx, nil
}
