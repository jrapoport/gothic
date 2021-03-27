package confirm

import (
	"errors"
	"net/http"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/hosts/rest"
)

// TODO: Merge this into the account endpoint to support authenticated
// 	account confirmation in lieu of an email address.
const (
	// Confirm is the confirmation rest endpoint.
	Confirm = "/confirm"
)

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
	r.Authenticated().Route(Confirm, func(rt *rest.Router) {
		rt.Post(rest.Root, s.SendConfirmUser)
	})
}

// SendConfirmUser resends a confirmation email to a user.
func (s *confirmServer) SendConfirmUser(w http.ResponseWriter, r *http.Request) {
	uid, err := rest.GetUserID(r)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	u, err := s.GetAuthenticatedUser(uid)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	ctx := rest.FromRequest(r)
	err = s.API.SendConfirmUser(ctx, u.ID)
	if errors.Is(err, config.ErrRateLimitExceeded) {
		s.ResponseCode(w, http.StatusTooEarly, err)
		return
	}
	s.ResponseError(w, err)
}
