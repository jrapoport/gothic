package users

import (
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/store"
)

// SearchRequest is an user search request
type SearchRequest struct {
	Filters store.Filters `json:"filters"  form:"filters"`
}

// NewSearchRequest returns a search request from a Request
func NewSearchRequest(r *http.Request) (*SearchRequest, error) {
	req := new(SearchRequest)
	data := store.Filters{}
	err := rest.UnmarshalRequest(r, &data)
	if err != nil {
		return nil, err
	}
	delete(data, key.Sort)
	delete(data, key.Page)
	delete(data, key.PageSize)
	req.Filters = data
	return req, nil
}

func (s *usersServer) SearchUsers(w http.ResponseWriter, r *http.Request) {
	req, err := NewSearchRequest(r)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	page := rest.PaginateRequest(r)
	ctx := rest.FromRequest(r)
	users, err := s.API.SearchUsers(ctx, req.Filters, page)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	s.PagedResponse(w, r, users, page)
}
