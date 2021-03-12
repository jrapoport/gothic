package signup

import (
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/store/types"
)

// Endpoint is the signup endpoint.
const Endpoint = "/signup"

// Request is an signup server request
type Request struct {
	Email    string    `json:"email" form:"email"`
	Username string    `json:"username" form:"username"`
	Password string    `json:"password" form:"password"`
	Data     types.Map `json:"data" form:"data"`
}

type signupServer struct {
	*rest.Server
}

func newSignupServer(srv *rest.Server) *signupServer {
	srv.FieldLogger = srv.WithField("module", "signup")
	return &signupServer{srv}
}

// RegisterServer registers a new signup server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newSignupServer(srv))
}

func register(s *http.Server, srv *signupServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *signupServer) addRoutes(r *rest.Router) {
	r.Post(Endpoint, s.Signup)
}

func (s *signupServer) Signup(w http.ResponseWriter, r *http.Request) {
	req := new(Request)
	err := rest.UnmarshalRequest(r, req)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	s.Debugf("signup user: %v", req)
	ctx := rest.FromRequest(r)
	u, err := s.API.Signup(ctx, req.Email, req.Username, req.Password, req.Data)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	res := rest.NewUserResponse(u)
	s.Debugf("signed up user: %v", res)
	if s.Config().MaskEmails {
		res.MaskEmail()
	}
	bt, err := s.GrantBearerToken(ctx, u)
	if err != nil {
		s.AuthError(w, err)
		return
	}
	res.Token = rest.NewBearerResponse(bt)
	s.AuthResponse(w, r, bt.String(), res)
}
