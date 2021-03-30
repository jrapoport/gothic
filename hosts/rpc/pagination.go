package rpc

import (
	api "github.com/jrapoport/gothic/api/grpc/rpc"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
)

// PaginateRequest paginates an http request
func PaginateRequest(req *api.SearchRequest) *store.Pagination {
	page := utils.Max(int(req.GetPage()), 1)
	perPage := store.MaxPageSize
	if req.PageSize != nil {
		perPage = int(req.GetPageSize())
		perPage = utils.Clamp(perPage, 1, store.MaxPageSize)
	}
	return &store.Pagination{
		Page: page,
		Size: perPage,
	}
}

func PaginateResponse(page *store.Pagination) *api.PagedResponse {
	res := &api.PagedResponse{
		Number: int64(page.Page),
		Count:  int64(page.Count),
		Total:  page.Total,
		Size:   int64(page.Length),
	}
	if page.Page >= 1 {
		res.First = 1
	}
	if page.Prev != 0 {
		res.Prev = int64(page.Prev)
	}
	if page.Next != 0 {
		res.Next = int64(page.Next)
	}
	if page.Page != page.Count {
		res.Last = int64(page.Count)
	}
	return res
}
