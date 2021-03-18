package users

import (
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/types"
)

// CreateRequest is an signup server request
type CreateRequest struct {
	Email    string    `json:"email" form:"email"`
	Username string    `json:"username" form:"username"`
	Password string    `json:"password" form:"password"`
	Data     types.Map `json:"data" form:"data"`
	Admin    bool      `json:"admin" form:"admin" `
}

// AdminCreateUser creates a user.
func (s *usersServer) AdminCreateUser(w http.ResponseWriter, r *http.Request) {
	uid, err := rest.GetUserID(r)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	req := new(CreateRequest)
	err = rest.UnmarshalRequest(r, req)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	s.Debugf("create user %s: %v", uid.String(), req)
	ctx := rest.FromRequest(r)
	u, err := s.API.AdminCreateUser(ctx, req.Email, req.Username, req.Password, req.Data, req.Admin)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	res := rest.NewUserResponse(u)
	s.Debugf("created user %s: %v", uid, res)
	s.Response(w, res)
}
