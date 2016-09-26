package api

import (
	"context"
	"net/http"
)

const description = `{
  "name": "Authlify",
  "description": "Authlify is a user registration and authentication API"
}`

// Index shows a description of the API
func (a *API) Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(description))
}
