package invite

import (
	"errors"
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
)

// Endpoint is the endpoint for an invite server.
const Endpoint = "/invite"

// Request is an invite server request
type Request struct {
	Email string `json:"email" form:"email"`
}

type inviteServer struct {
	*rest.Server
}

func newInviteServer(srv *rest.Server) *inviteServer {
	srv.FieldLogger = srv.WithField("module", "config")
	return &inviteServer{srv}
}

// RegisterServer an invite server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newInviteServer(srv))
}

func register(s *http.Server, srv *inviteServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *inviteServer) addRoutes(r *rest.Router) {
	r.Authenticated().Confirmed().Post(Endpoint, s.SendInviteUser)
}

// SendInviteUser sends an invite to an email address.
func (s *inviteServer) SendInviteUser(w http.ResponseWriter, r *http.Request) {
	uid, err := rest.GetUserID(r)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	req := new(Request)
	err = rest.UnmarshalRequest(r, req)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	if req.Email == "" {
		err = errors.New("email not found")
		s.ResponseCode(w, http.StatusUnprocessableEntity, err)
		return
	}
	ctx := rest.FromRequest(r)
	err = s.API.SendInviteUser(ctx, uid, req.Email)
	s.Response(w, err)
}
