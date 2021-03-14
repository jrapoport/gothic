package confirm

import (
	"errors"
	"github.com/jrapoport/gothic/config"
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
)

const (
	// Endpoint is the confirmation rest endpoint.
	Endpoint = "/confirm"
	// Root is the base route.
	Root = "/"
	// Send is the route to resent a confirmation email.
	Send = "/send"
)

// Request is an confirm server request
type Request struct {
	Email string `json:"email" form:"email"`
	Token string `json:"token" form:"token"`
}

type confirmServer struct {
	*rest.Server
}

func newConfirmServer(srv *rest.Server) *confirmServer {
	srv.FieldLogger = srv.WithField("module", "confirm")
	return &confirmServer{srv}
}

// RegisterServer registers a confirmation server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newConfirmServer(srv))
}

func register(s *http.Server, srv *confirmServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *confirmServer) addRoutes(r *rest.Router) {
	r.Route(Endpoint, func(rt *rest.Router) {
		rt.Post(Root, s.ConfirmUser)
		rt.Post(Send, s.SendConfirmUser)
	})
}

// ConfirmUser confirms a user with a token.
func (s *confirmServer) ConfirmUser(w http.ResponseWriter, r *http.Request) {
	req := new(Request)
	err := rest.UnmarshalRequest(r, req)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	if req.Token == "" {
		s.ResponseCode(w, http.StatusUnprocessableEntity, err)
		return
	}
	s.Debugf("confirm user: %v", req)
	ctx := rest.FromRequest(r)
	u, err := s.API.ConfirmUser(ctx, req.Token)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	bt, err := s.GrantBearerToken(ctx, u)
	if err != nil {
		s.AuthError(w, err)
		return
	}
	s.Debugf("confirmed user: %s", bt.UserID)
	s.AuthResponse(w, r, bt.Token, bt)
}

// SendConfirmUser resends a confirmation email to a user.
func (s *confirmServer) SendConfirmUser(w http.ResponseWriter, r *http.Request) {
	req := new(Request)
	err := rest.UnmarshalRequest(r, req)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	if req.Email == "" {
		err = errors.New("email not found")
		s.ResponseCode(w, http.StatusUnprocessableEntity, err)
		return
	}
	email, err := s.ValidateEmail(req.Email)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	u, err := s.GetUserWithEmail(email)
	if err != nil {
		s.AuthError(w, err)
		return
	}
	ctx := rest.FromRequest(r)
	err = s.API.SendConfirmUser(ctx, u.ID)
	if errors.Is(err, config.ErrRateLimitExceeded) {
		s.ResponseCode(w, http.StatusTooEarly, err)
		return
	}
	s.AuthError(w, err)
}
