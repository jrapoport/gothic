package codes

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/code"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/types/key"
)

// Codes endpoint
const (
	Codes  = "/codes"
	Create = rest.Root
	Read   = "/{" + key.Code + "}"
	Delete = Read
)

type codesServer struct {
	*rest.Server
}

func newSignupServer(srv *rest.Server) *codesServer {
	srv.Logger = srv.WithName("user")
	return &codesServer{srv}
}

// RegisterServer registers a new user server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newSignupServer(srv))
}

func register(s *http.Server, srv *codesServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *codesServer) addRoutes(r *rest.Router) {
	r.Authenticated().Admin().Route(Codes, func(rt *rest.Router) {
		rt.Post(Create, s.CreateSignupCodes)
		rt.Get(Read, s.CheckSignupCode)
		rt.Delete(Delete, s.DeleteSignupCode)
	})
}

// CreateSignupCodes creates a list of signup codes.
func (s *codesServer) CreateSignupCodes(w http.ResponseWriter, r *http.Request) {
	// CreateRequest is an user server request
	type Request struct {
		Uses  int `json:"uses" form:"uses"`
		Count int `json:"count" form:"count"`
	}
	req := new(Request)
	err := rest.UnmarshalRequest(r, req)
	if err != nil {
		s.ResponseCode(w, http.StatusUnprocessableEntity, err)
		return
	}
	_, err = s.ValidateAdmin(r)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	ctx := rest.FromRequest(r)
	s.Debugf("create signup codes %s", ctx.AdminID())
	list, err := s.API.CreateSignupCodes(ctx, req.Uses, req.Count)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	s.Debugf("created signup codes %v", list)
	s.Response(w, list)
}

// CheckSignupCode checks a code
func (s *codesServer) CheckSignupCode(w http.ResponseWriter, r *http.Request) {
	tok := rest.URLParam(r, key.Code)
	if tok == "" {
		err := errors.New("invalid code")
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	// Response contains an http signup code response.
	type Response struct {
		Valid      bool          `json:"valid"`
		Code       string        `json:"code,omitempty"`
		Format     code.Format   `json:"format"`
		Type       token.Type    `json:"type"`
		Expiration time.Duration `json:"expiration,omitempty"`
		UserID     uuid.UUID     `json:"user_id"`
	}
	s.Debugf("check signup code: %s", tok)
	sc, err := s.API.CheckSignupCode(tok)
	if err != nil && errors.Is(err, code.ErrUnusableCode) {
		s.Debugf("checked signup code is invalid: %s", tok)
		res := &Response{
			Valid: false,
			Code:  tok,
		}
		s.Response(w, res)
		return
	} else if err != nil {
		s.ResponseCode(w, http.StatusInternalServerError, err)
		return
	}
	s.Debugf("checked signup code is valid: %s", tok)
	res := &Response{
		Valid:      true,
		Code:       sc.Token,
		Format:     sc.Format,
		Type:       sc.Type,
		Expiration: sc.Expiration,
		UserID:     sc.UserID,
	}
	s.Response(w, res)
}

// CheckSignupCode checks a code
func (s *codesServer) DeleteSignupCode(w http.ResponseWriter, r *http.Request) {
	tok := rest.URLParam(r, key.Code)
	if tok == "" {
		err := errors.New("invalid code")
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	s.Debugf("delete signup code: %s", tok)
	err := s.API.DeleteSignupCode(tok)
	if err != nil {
		s.ResponseCode(w, http.StatusInternalServerError, err)
		return
	}
	s.Response(w, nil)
}
