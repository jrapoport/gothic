package rest

import (
	"net/http"

	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/store"
	"github.com/segmentio/encoding/json"
)

// Server represents an REST server.
type Server struct {
	*core.Server
}

// NewServer creates a new REST Server.
func NewServer(s *core.Server) *Server {
	return &Server{s}
}

// Clone returns a clone of the server.
func (s *Server) Clone() *Server {
	return &Server{s.Server.Clone()}
}

// Response wraps an http JSONContent response. If v is an
// standard error or Error, it writes an Error instead.
func (s *Server) Response(w http.ResponseWriter, v interface{}) {
	if v == nil {
		s.ResponseCode(w, http.StatusOK, nil)
		return
	}
	switch val := v.(type) {
	case error:
		s.ResponseCode(w, http.StatusInternalServerError, val)
		return
	case *tokens.BearerToken:
		v = NewBearerResponse(val)
		s.Debugf("returned bearer token: %v", v)
		break
	default:
		break
	}
	b, err := json.Marshal(v)
	if err != nil {
		s.ResponseCode(w, http.StatusInternalServerError, err)
		return
	}
	s.Debugf("response: %s", string(b))
	w.Header().Set(ContentType, JSONContent)
	// ResponseWriter.Write() calls w.WriteHeader(http.StatusOK)
	if _, err = w.Write(b); err != nil {
		s.ResponseCode(w, http.StatusInternalServerError, err)
		return
	}
}

// ResponseCode logs an error and the writes an sanitized standard response.
func (s *Server) ResponseCode(w http.ResponseWriter, code int, err error) {
	if err != nil {
		s.Error(err)
	}
	ResponseCode(w, code, err)
}

// ResponseError logs an error and the writes an sanitized standard response.
func (s *Server) ResponseError(w http.ResponseWriter, err error) {
	code := http.StatusOK
	if err != nil {
		code = http.StatusInternalServerError
	}
	s.ResponseCode(w, code, err)
}

// AuthResponse will log the error but hide it so we don't leak information
func (s *Server) AuthResponse(w http.ResponseWriter, r *http.Request, tok string, v interface{}) {
	UseCookie(w, r, tok, s.Config().Cookies.Duration)
	s.Response(w, v)
}

// AuthError will log the error but hide it so we don't leak information
func (s *Server) AuthError(w http.ResponseWriter, err error) {
	ClearCookie(w)
	s.ResponseCode(w, http.StatusOK, err)
}

// PagedResponse will return a pages response.
func (s *Server) PagedResponse(w http.ResponseWriter, r *http.Request,
	v interface{}, page *store.Pagination) {
	PaginateResponse(w, r, page)
	s.Response(w, v)
}

// ResponseCode writes a standard http response
func ResponseCode(w http.ResponseWriter, code int, err error) {
	if code == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		return
	}
	msg := http.StatusText(code)
	if err != nil && err.Error() == "" {
		msg = err.Error()
	}
	http.Error(w, msg, code)
}
