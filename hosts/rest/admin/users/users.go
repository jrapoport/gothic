package users

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
)

// Users endpoint
const (
	Users    = "/users"
	Search   = rest.Root
	UserID   = "/{" + key.UserID + "}" // select a user
	Create   = rest.Root
	Read     = rest.Root
	Update   = rest.Root
	Delete   = rest.Root
	Promote  = "/promote"
	Metadata = "/metadata"
)

// Request is an user server request
type Request struct {
	Username    string    `json:"username,omitempty" form:"username"`
	Data        types.Map `json:"data,omitempty" form:"data"`
	Metadata    types.Map `json:"metadata,omitempty" form:"metadata"`
	Password    string    `json:"password,omitempty" form:"password"`
	OldPassword string    `json:"new_password,omitempty" form:"new_password"`
}

type usersServer struct {
	*rest.Server
}

func newUserServer(srv *rest.Server) *usersServer {
	srv.Logger = srv.WithName("user")
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
			uid.Put(Update, s.AdminUpdateUser)
			uid.Delete(Delete, s.AdminDeleteUser)
			uid.Post(Promote, s.AdminPromoteUser)
			uid.Post(Metadata, s.AdminUpdateUserMetadata)
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
	_, err = s.ValidateAdmin(r)
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

// AdminUpdateUser updates a user.
func (s *usersServer) AdminUpdateUser(w http.ResponseWriter, r *http.Request) {
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
	_, err = s.ValidateAdmin(r)
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
	_, err = s.ValidateAdmin(r)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	hard := rest.URLParam(r, key.Hard) == "true"
	ctx := rest.FromRequest(r)
	s.Debugf("delete user %", uid.String())
	err = s.API.DeleteUser(ctx, uid, hard)
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
	_, err = s.ValidateAdmin(r)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	ctx := rest.FromRequest(r)
	s.Debugf("promote user %", uid.String())
	u, err := s.API.PromoteUser(ctx, uid)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	res := rest.NewUserResponse(u)
	s.Debugf("promoted user %s: %v", uid, res)
	s.Response(w, res)
}

// AdminUpdateUserMetadata updates a user's metadata.
func (s *usersServer) AdminUpdateUserMetadata(w http.ResponseWriter, r *http.Request) {
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
	_, err = s.ValidateAdmin(r)
	if err != nil {
		s.ResponseCode(w, http.StatusUnauthorized, err)
		return
	}
	ctx := rest.FromRequest(r)
	s.Debugf("update user %s: %v", uid.String(), req)
	u, err := s.API.UpdateUserMetadata(ctx, uid, req.Metadata)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	s.Debugf("updated user metadata %s: %v", uid, u.Metadata)
	s.Response(w, u.Metadata)
}
