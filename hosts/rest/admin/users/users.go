package users

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
)

// Users endpoint
const (
	Users   = "/users"
	Search  = rest.Root
	UserID  = "/{" + key.UserID + "}" // select a user
	Create  = rest.Root
	Read    = rest.Root
	Update  = rest.Root
	Delete  = rest.Root
	Promote = "/promote"
)

// Request is an user server request
type Request struct {
	Username    string    `json:"username" form:"username"`
	Data        types.Map `json:"data" form:"data"`
	Password    string    `json:"password" form:"password"`
	OldPassword string    `json:"new_password" form:"new_password"`
}

type usersServer struct {
	*rest.Server
}

func newUserServer(srv *rest.Server) *usersServer {
	srv.FieldLogger = srv.WithField("service", "user")
	return &usersServer{srv}
}

// RegisterServer registers a new user server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newUserServer(srv))
}

func register(s *http.Server, srv *usersServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *usersServer) addRoutes(r *rest.Router) {
	r.Authenticated().Admin().Route(Users, func(rt *rest.Router) {
		rt.Get(Search, s.SearchUsers)
		rt.Post(Create, s.AdminCreateUser)
		rt.Route(UserID, func(uid *rest.Router) {
			uid.Get(Read, s.GetUser)
			uid.Put(Update, s.UpdateUser)
			uid.Delete(Delete, s.AdminDeleteUser)
			uid.Post(Promote, s.AdminPromoteUser)
		})
	})
}

// GetUser gets a user.
func (s *usersServer) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := rest.URLParam(r, key.UserID)
	uid, err := uuid.Parse(userID)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	err = s.validateAdmin(r)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	s.Debugf("get user %s", uid)
	u, err := s.API.GetUser(uid)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	res := rest.NewUserResponse(u)
	s.Debugf("got user %s: %v", uid, res)
	s.Response(w, res)
}

// UpdateUser updates a user.
func (s *usersServer) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := rest.URLParam(r, key.UserID)
	uid, err := uuid.Parse(userID)
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
	err = s.validateAdmin(r)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	ctx := rest.FromRequest(r)
	s.Debugf("update user %s: %v", uid.String(), req)
	u, err := s.API.UpdateUser(ctx, uid, &req.Username, req.Data)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	res := rest.NewUserResponse(u)
	s.Debugf("updated user %s: %v", uid, res)
	s.Response(w, res)
}

// AdminDeleteUser updates a user.
func (s *usersServer) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := rest.URLParam(r, key.UserID)
	uid, err := uuid.Parse(userID)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	err = s.validateAdmin(r)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	ctx := rest.FromRequest(r)
	s.Debugf("delete user %", uid.String())
	err = s.API.AdminDeleteUser(ctx, uid)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	s.Debugf("deleted user %s", uid)
	s.Response(w, nil)
}

// AdminPromoteUser updates a user.
func (s *usersServer) AdminPromoteUser(w http.ResponseWriter, r *http.Request) {
	userID := rest.URLParam(r, key.UserID)
	uid, err := uuid.Parse(userID)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	err = s.validateAdmin(r)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	ctx := rest.FromRequest(r)
	s.Debugf("promote user %", uid.String())
	u, err := s.API.AdminPromoteUser(ctx, uid)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	res := rest.NewUserResponse(u)
	s.Debugf("promoted user %s: %v", uid, res)
	s.Response(w, res)
}

func (s *usersServer) validateAdmin(r *http.Request) error {
	aid, err := rest.GetUserID(r)
	if err != nil {
		return err
	}
	adm, err := s.GetAuthenticatedUser(aid)
	if err != nil {
		return err
	}
	if !adm.IsAdmin() {
		err = fmt.Errorf("admin required: %s", adm.ID)
		return err
	}
	return nil
}
