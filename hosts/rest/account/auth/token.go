package auth

import (
	"errors"
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
)

func (s *authServer) RefreshBearerToken(w http.ResponseWriter, r *http.Request) {
	req := new(Request)
	err := rest.UnmarshalRequest(r, req)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	if req.Token == "" {
		err = errors.New("refresh token not found")
		s.ResponseCode(w, http.StatusUnprocessableEntity, err)
		return
	}
	ctx := rest.FromRequest(r)
	s.Debugf("refresh token: %v", req)
	bt, err := s.API.RefreshBearerToken(ctx, req.Token)
	if err != nil {
		s.AuthError(w, err)
		return
	}
	s.AuthResponse(w, r, bt.Token, bt)
}
