package login

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/hosts/rest"
)

func (s *loginServer) Logout(w http.ResponseWriter, r *http.Request) {
	rtx := rest.FromRequest(r)
	uid := rtx.UserID()
	if uid == uuid.Nil {
		err := errors.New("invalid user id")
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	s.Debugf("logout user: %s (%v)", uid, rtx.Provider())
	rest.ClearCookie(w)
	err := s.API.Logout(rtx, uid)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	s.ResponseCode(w, http.StatusOK, nil)
}
