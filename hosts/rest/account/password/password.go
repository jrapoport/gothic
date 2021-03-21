package password

import (
	"errors"
	"net/http"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/hosts/rest"
)

const (
	// Endpoint for password
	Endpoint = "/password"
	// Confirm confirms a password
	Confirm = "/"
	// Reset starts a password reset
	Reset = "/reset"
)

// Request is an password server request
type Request struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Token    string `json:"token" form:"token"`
}

type passwordServer struct {
	*rest.Server
}

func newPasswordServer(srv *rest.Server) *passwordServer {
	srv.FieldLogger = srv.WithField("module", "password")
	return &passwordServer{srv}
}

// RegisterServer registers a new password server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newPasswordServer(srv))
}

func register(s *http.Server, srv *passwordServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *passwordServer) addRoutes(r *rest.Router) {
	r.Route(Endpoint, func(rt *rest.Router) {
		rt.Post(Reset, s.SendResetPassword)
		rt.Post(Confirm, s.ConfirmResetPassword)
	})
}

func (s *passwordServer) SendResetPassword(w http.ResponseWriter, r *http.Request) {
	req := new(Request)
	err := rest.UnmarshalRequest(r, req)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	if req.Email == "" {
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
	err = s.API.SendResetPassword(ctx, u.ID)
	if errors.Is(err, config.ErrRateLimitExceeded) {
		s.ResponseCode(w, http.StatusTooEarly, err)
		return
	}
	// log any error but hide it so we don't leak information
	s.AuthError(w, err)
}

func (s *passwordServer) ConfirmResetPassword(w http.ResponseWriter, r *http.Request) {
	req := new(Request)
	err := rest.UnmarshalRequest(r, req)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	if req.Password == "" {
		err = errors.New("password not found")
		s.ResponseCode(w, http.StatusUnprocessableEntity, err)
		return
	}
	if req.Token == "" {
		err = errors.New("token not found")
		s.ResponseCode(w, http.StatusUnprocessableEntity, err)
		return
	}
	s.Debugf("change password: %v", req)
	ctx := rest.FromRequest(r)
	u, err := s.API.ConfirmResetPassword(ctx, req.Token, req.Password)
	if err != nil {
		s.AuthError(w, err)
		return
	}
	bt, err := s.GrantBearerToken(ctx, u)
	if err != nil {
		s.AuthError(w, err)
		return
	}
	s.Debugf("password changed: %s", bt.UserID)
	s.AuthResponse(w, r, bt.Token, bt)
}
