package login

import (
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
)

// Login endpoints
const (
	Login  = "/login"
	Logout = "/logout"
)

// Request is an login server request
type Request struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type loginServer struct {
	*rest.Server
}

func newLoginServer(srv *rest.Server) *loginServer {
	srv.Logger = srv.WithName("login")
	return &loginServer{srv}
}

// RegisterServer registers a new login server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newLoginServer(srv))
}

func register(s *http.Server, srv *loginServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *loginServer) addRoutes(r *rest.Router) {
	r.Post(Login, s.Login)
	r.Authenticated().Get(Logout, s.Logout)
}

func (s *loginServer) Login(w http.ResponseWriter, r *http.Request) {
	req := new(Request)
	err := rest.UnmarshalRequest(r, req)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	s.Debugf("login user: %v (%v)", req)
	ctx := rest.FromRequest(r)
	u, err := s.API.Login(ctx, req.Email, req.Password)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	res := rest.NewUserResponse(u)
	if s.Config().MaskEmails {
		res.MaskEmail()
	}
	bt, err := s.GrantBearerToken(ctx, u)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	res.Token = rest.NewBearerResponse(bt)
	s.Debugf("logged in user: %s", bt.UserID)
	s.AuthResponse(w, r, bt.String(), res)
}
