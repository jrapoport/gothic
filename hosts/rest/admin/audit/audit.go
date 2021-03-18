package audit

import (
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/store"
)

// Endpoint is the config endpoint
const Endpoint = "/audit"

// Request is an audit log search request
type Request struct {
	Filters store.Filters `json:"filters"  form:"filters"`
}

// NewRequest returns a search request from a Request
func NewRequest(r *http.Request) (*Request, error) {
	req := new(Request)
	data := store.Filters{}
	err := rest.UnmarshalRequest(r, &data)
	if err != nil {
		return nil, err
	}
	delete(data, key.Sort)
	delete(data, key.Page)
	delete(data, key.PageCount)
	req.Filters = data
	return req, nil
}

type auditServer struct {
	*rest.Server
}

// NewAuditServer returns a new config rest server.
func newAuditServer(srv *rest.Server) *auditServer {
	srv.FieldLogger = srv.WithField("module", "audit")
	return &auditServer{srv}
}

// RegisterServer registers a new audit server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newAuditServer(srv))
}

func register(s *http.Server, srv *auditServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *auditServer) addRoutes(r *rest.Router) {
	r.Get(Endpoint, s.SearchAuditLogs)
}

func (s *auditServer) SearchAuditLogs(w http.ResponseWriter, r *http.Request) {
	req, err := NewRequest(r)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	page, err := rest.PaginateRequest(r)
	if err != nil {
		s.ResponseCode(w, http.StatusBadRequest, err)
		return
	}
	ctx := rest.FromRequest(r)
	logs, err := s.API.SearchAuditLogs(ctx, req.Filters, page)
	if err != nil {
		s.ResponseError(w, err)
		return
	}
	s.PagedResponse(w, r, logs, page)
}
