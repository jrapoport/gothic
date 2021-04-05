package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
)

const (
	// Link http header.
	Link = "Link"
	// PageNumber is the current page number.
	PageNumber = "X-Page-Number"
	// PageCount is the total pages of pages.
	PageCount = "X-Page-Count"
	// PageSize is the number of items in the page.
	PageSize = "X-Page-Size"
	// PageTotal is the total number of items across all pages.
	PageTotal = "X-Page-Total"
)

// PaginateRequest paginates an http request
func PaginateRequest(r *http.Request) *store.Pagination {
	if r.Form == nil {
		const defaultMaxMemory = 32 << 20 // 32 MB
		// it looks like it will return "no body" but still do the right
		// thing with the form. net/http/request.go ignores it also.
		_ = r.ParseMultipartForm(defaultMaxMemory)
	}
	var page = 1
	if v := r.FormValue(key.Page); v != "" {
		num, _ := strconv.ParseInt(v, 10, 64)
		page = utils.Max(int(num), page)
	}
	var perPage = store.MaxPageSize
	if v := r.FormValue(key.PageSize); v != "" {
		per, _ := strconv.ParseInt(v, 10, 64)
		perPage = utils.Clamp(int(per), 1, perPage)
	}
	return &store.Pagination{
		Index: page,
		Size:  perPage,
	}
}

const (
	firstRel = "first"
	prevRel  = "prev"
	nextRel  = "next"
	lastRel  = "last"
)

func linkRel(u *url.URL, rel string) string {
	return fmt.Sprintf(`<%s>; rel="%s"`, u.String(), rel)
}

func firstLink(u *url.URL, page *store.Pagination) string {
	if page.Index <= 1 {
		return ""
	}
	u.RawQuery = url.Values{
		key.Page:    []string{"1"},
		key.PerPage: []string{strconv.Itoa(page.Size)},
	}.Encode()
	return linkRel(u, firstRel)
}

func prevLink(u *url.URL, page *store.Pagination) string {
	if page.Prev == 0 {
		return ""
	}
	u.RawQuery = url.Values{
		key.Page:    []string{strconv.Itoa(page.Prev)},
		key.PerPage: []string{strconv.Itoa(page.Size)},
	}.Encode()
	return linkRel(u, prevRel)
}

func nextLink(u *url.URL, page *store.Pagination) string {
	if page.Next == 0 {
		return ""
	}
	u.RawQuery = url.Values{
		key.Page:    []string{strconv.Itoa(page.Next)},
		key.PerPage: []string{strconv.Itoa(page.Size)},
	}.Encode()
	return linkRel(u, nextRel)
}

func lastLink(u *url.URL, page *store.Pagination) string {
	if page.Index == page.Count {
		return ""
	}
	u.RawQuery = url.Values{
		key.Page:    []string{strconv.Itoa(page.Count)},
		key.PerPage: []string{strconv.Itoa(page.Size)},
	}.Encode()
	return linkRel(u, lastRel)
}

// PaginateResponse writes a paginated http response.
func PaginateResponse(w http.ResponseWriter, r *http.Request, page *store.Pagination) {
	w.Header().Add(PageNumber, strconv.Itoa(page.Index))
	w.Header().Add(PageCount, strconv.Itoa(page.Count))
	w.Header().Add(PageTotal, strconv.FormatUint(page.Total, 10))
	w.Header().Add(PageSize, strconv.Itoa(page.Length))
	u := &url.URL{}
	u.Scheme = "http"
	if r.TLS != nil {
		u.Scheme = "https"
	}
	if s := r.Header.Get(ForwardedProto); s != "" {
		u.Scheme = s
	}
	u.Host = r.Host
	u.Path = r.URL.Path
	var links []string
	if l := firstLink(u, page); l != "" {
		links = append(links, l)
	}
	if l := prevLink(u, page); l != "" {
		links = append(links, l)
	}
	if l := nextLink(u, page); l != "" {
		links = append(links, l)
	}
	if l := lastLink(u, page); l != "" {
		links = append(links, l)
	}
	linkHeader := strings.Join(links, ",")
	w.Header().Add(Link, linkHeader)
}
