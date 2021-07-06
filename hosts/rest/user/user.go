package user

import (
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/modules/invite"
	"github.com/jrapoport/gothic/hosts/rest/user/email"
	"github.com/jrapoport/gothic/store/types"
)

const (
	// Endpoint is the user endpoint.
	Endpoint = "/user"
	// Root gets or updates a user.
	Root = "/"
	// Password changes a user password.
	Password = "/password"
)

type Request struct {
	Username    string    `json:"username" form:"username"`
	Data        types.Map `json:"data" form:"data"`
	Password    string    `json:"password" form:"password"`
	OldPassword string    `json:"old_password" form:"old_password"`
}

type userServer struct {
	*rest.Server
}

func newUserServer(srv *rest.Server) *userServer {
	srv.FieldLogger = srv.WithField("service", "user")
	return &userServer{srv}
}

// RegisterServer registers a new user server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newUserServer(srv))
}

func register(s *http.Server, srv *userServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *userServer) addRoutes(r *rest.Router) {
	r.Authenticated().Confirmed().Route(Endpoint, func(rt *rest.Router) {
		rt.Get(Root, s.GetUser)
		rt.Put(Root, s.UpdateUser)
		rt.Put(Password, s.ChangePassword)
		email.RegisterServer(&http.Server{Handler: rt}, s.Clone())
		invite.RegisterServer(&http.Server{Handler: rt}, s.Clone())
	})
}

// GetUser gets a user.
func (s *userServer) GetUser(w http.ResponseWriter, r *http.Request) {
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
	s.Debugf("get user %s", u.ID)
	res := rest.NewUserResponse(u)
	if s.Config().MaskEmails {
		res.MaskEmail()
	}
	s.Debugf("got user %s: %v", uid, res)
	s.Response(w, res)
}

// UpdateUser updates a user.
func (s *userServer) UpdateUser(w http.ResponseWriter, r *http.Request) {
	uid, err := rest.GetUserID(r)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	req := new(Request)
	err = rest.UnmarshalRequest(r, req)
	if err != nil {
		s.ResponseCode(w, http.StatusUnprocessableEntity, err)
		return
	}
	_, err = s.GetAuthenticatedUser(uid)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	s.Debugf("update user %s: %v", uid.String(), req)
	ctx := rest.FromRequest(r)
	u, err := s.API.UpdateUser(ctx, uid, &req.Username, req.Data)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	res := rest.NewUserResponse(u)
	if s.Config().MaskEmails {
		res.MaskEmail()
	}
	s.Debugf("updated user %s: %v", uid, res)
	s.Response(w, res)
}

func (s *userServer) ChangePassword(w http.ResponseWriter, r *http.Request) {
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
	old := req.OldPassword
	u, err := s.GetAuthenticatedUser(uid)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	s.Debugf("change password %s: %v", u.ID, req)
	ctx := rest.FromRequest(r)
	u, err = s.API.ChangePassword(ctx, u.ID, old, req.Password)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
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
