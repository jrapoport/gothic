package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/utils"
)

// Auth endpoints
const (
	Auth = "/auth"
	// Refresh refreshes a user refresh token
	Refresh = rest.Root
	// Provider returns the auth url for a provider.
	Provider = "/{" + key.Provider + "}" //
	// Callback authorizes a user.
	Callback = "/callback"
)

// Request is an auth server request
type Request struct {
	State string `json:"state" form:"state"`
	Token string `json:"token" form:"token"`
}

type authServer struct {
	*rest.Server
}

func newAuthServer(srv *rest.Server) *authServer {
	srv.FieldLogger = srv.WithField("module", "authorize")
	return &authServer{srv}
}

// RegisterServer an auth rest server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newAuthServer(srv))
}

func register(s *http.Server, srv *authServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *authServer) addRoutes(r *rest.Router) {
	r.Route(Auth, func(rt *rest.Router) {
		rt.Post(Refresh, s.RefreshBearerToken)
		rt.Get(Provider, s.GetAuthorizationURL)
		rt.Post(Callback, s.AuthorizeUser)
	})
}

// GetAuthorizationURL returns an auth url for a provider.
func (s *authServer) GetAuthorizationURL(w http.ResponseWriter, r *http.Request) {
	p := provider.Name(rest.URLParam(r, key.Provider))
	if !p.IsExternal() {
		err := fmt.Errorf("invalid provider: %s", p)
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	s.Debugf("get authorization url for %s: (%v)", p)
	ctx := rest.FromRequest(r)
	ctx.SetProvider(p)
	au, err := s.API.GetAuthorizationURL(ctx, p)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	s.Debugf("got authorization url: %s", au)
	http.Redirect(w, r, au, http.StatusFound)
}

func (s *authServer) AuthorizeUser(w http.ResponseWriter, r *http.Request) {
	req := new(Request)
	err := rest.UnmarshalRequest(r, req)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	if req.State == "" {
		err = errors.New("state not found")
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	data := utils.URLValuesToMap(r.Form, false)
	ctx := rest.FromRequest(r)
	u, err := s.API.AuthorizeUser(ctx, req.State, data)
	if err != nil {
		s.AuthError(w, err)
		return
	}
	bt, err := s.GrantBearerToken(ctx, u)
	if err != nil {
		s.AuthError(w, err)
		return
	}
	res := rest.NewUserResponse(u)
	if s.Config().MaskEmails {
		res.MaskEmail()
	}
	res.Token = rest.NewBearerResponse(bt)
	s.Debugf("authorized user: %v", res)
	s.AuthResponse(w, r, bt.String(), res)
}
