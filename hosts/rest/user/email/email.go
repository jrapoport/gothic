package email

import (
	"errors"
	"net/http"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/jwt"
	"github.com/jrapoport/gothic/models/types/key"
)

const (
	// Email for the email changes.
	Email = "/email"
	// Change is the email change endpoint.
	Change = "/change"
	// Confirm an email change.
	Confirm = "/confirm"
)

// Request is an email server request
type Request struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Token    string `json:"token" form:"token"`
}

type emailServer struct {
	*rest.Server
}

func newEmailServer(srv *rest.Server) *emailServer {
	srv.Logger = srv.WithName("email")
	return &emailServer{srv}
}

// RegisterServer registers a new email server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newEmailServer(srv))
}

func register(s *http.Server, srv *emailServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *emailServer) addRoutes(r *rest.Router) {
	r.Route(Email, func(rt *rest.Router) {
		rt.Post(Confirm, s.ConfirmChangeEmail)
		rt.Authenticated().Confirmed().Route(rest.Root, func(cr *rest.Router) {
			// cr.Post(rest.Root, s.UnMaskEmail)
			cr.Post(Change, s.SendChangeEmail)
		})
	})
}

func (s *emailServer) SendChangeEmail(w http.ResponseWriter, r *http.Request) {
	// we can safely ignore this error since this route is
	// protected we've already checked for a valid user id
	uid, _ := rest.GetUserID(r)
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
	u, err := s.GetAuthenticatedUser(uid)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	err = u.Authenticate(req.Password)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	ctx := rest.FromRequest(r)
	err = s.API.SendChangeEmail(ctx, u.ID, email)
	if errors.Is(err, config.ErrRateLimitExceeded) {
		s.ResponseCode(w, http.StatusTooEarly, err)
		return
	}
	s.AuthError(w, err)
}

func (s *emailServer) ConfirmChangeEmail(w http.ResponseWriter, r *http.Request) {
	req := new(Request)
	err := rest.UnmarshalRequest(r, req)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	if req.Token == "" {
		err = errors.New("token not found")
		s.ResponseCode(w, http.StatusUnprocessableEntity, err)
		return
	}
	data, err := jwt.ParseData(s.Config().JWT, req.Token)
	if err != nil {
		s.ResponseCode(w, http.StatusUnprocessableEntity, err)
		return
	}
	ct, ok := data[key.Token].(string)
	if !ok || ct == "" {
		err = errors.New("confirmation token not found")
		s.ResponseCode(w, http.StatusUnprocessableEntity, err)
		return
	}
	email, ok := data[key.Email].(string)
	if !ok || email == "" {
		err = errors.New("email not found")
		s.ResponseCode(w, http.StatusUnprocessableEntity, err)
		return
	}
	s.Debugf("change password: %v", req)
	ctx := rest.FromRequest(r)
	u, err := s.API.ConfirmChangeEmail(ctx, ct, email)
	if err != nil {
		s.ResponseError(w, err)
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

/*
func (s *emailServer) UnMaskEmail(w http.ResponseWriter, _ *http.Request) {
	s.ResponseCode(w, http.StatusOK, nil)
}
*/
