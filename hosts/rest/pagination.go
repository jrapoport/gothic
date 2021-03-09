package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types/key"
)

const (
	// Link http header.
	Link = "Link"
	// PageNumber is the current page number.
	PageNumber = "X-Page-Number"
	// PageCount is the total pages of pages.
	PageCount = "X-Page-Count"
	// PageLength is the actual number of items in the page.
	PageLength = "X-Page-Length"
	// PageTotal is the total number of items across all pages.
	PageTotal = "X-Page-Total"
)

// PaginateRequest paginates an http request
func PaginateRequest(r *http.Request) (*store.Pagination, error) {
	if r.Form == nil {
		const defaultMaxMemory = 32 << 20 // 32 MB
		// it looks like it will return "no body" but still do the right
		// thing with the form. net/http/request.go ignores it also.
		_ = r.ParseMultipartForm(defaultMaxMemory)
	}
	var page int64 = 1
	var perPage int64 = store.MaxPerPage
	var err error
	if v := r.FormValue(key.Page); v != "" {
		page, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if v := r.FormValue(key.PageCount); v != "" {
		perPage, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return &store.Pagination{
		Page: int(page),
		Size: int(perPage),
	}, nil
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

func firstLink(u *url.URL, p *store.Pagination) string {
	if p.Page <= 1 {
		return ""
	}
	u.RawQuery = url.Values{
		key.Page:    []string{"1"},
		key.PerPage: []string{strconv.Itoa(p.Size)},
	}.Encode()
	return linkRel(u, firstRel)
}

func prevLink(u *url.URL, p *store.Pagination) string {
	if p.Prev == 0 {
		return ""
	}
	u.RawQuery = url.Values{
		key.Page:    []string{strconv.Itoa(p.Prev)},
		key.PerPage: []string{strconv.Itoa(p.Size)},
	}.Encode()
	return linkRel(u, prevRel)
}

func nextLink(u *url.URL, p *store.Pagination) string {
	if p.Next == 0 {
		return ""
	}
	u.RawQuery = url.Values{
		key.Page:    []string{strconv.Itoa(p.Next)},
		key.PerPage: []string{strconv.Itoa(p.Size)},
	}.Encode()
	return linkRel(u, nextRel)
}

func lastLink(u *url.URL, p *store.Pagination) string {
	if p.Page == p.Count {
		return ""
	}
	u.RawQuery = url.Values{
		key.Page:    []string{strconv.Itoa(p.Count)},
		key.PerPage: []string{strconv.Itoa(p.Size)},
	}.Encode()
	return linkRel(u, lastRel)
}

// PaginateResponse writes a paginated http response.
func PaginateResponse(w http.ResponseWriter, r *http.Request, p *store.Pagination) {
	w.Header().Add(PageNumber, strconv.Itoa(p.Page))
	w.Header().Add(PageCount, strconv.Itoa(p.Count))
	w.Header().Add(PageTotal, strconv.FormatInt(p.Total, 10))
	w.Header().Add(PageLength, strconv.Itoa(p.Length))
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
	if l := firstLink(u, p); l != "" {
		links = append(links, l)
	}
	if l := prevLink(u, p); l != "" {
		links = append(links, l)
	}
	if l := nextLink(u, p); l != "" {
		links = append(links, l)
	}
	if l := lastLink(u, p); l != "" {
		links = append(links, l)
	}
	linkHeader := strings.Join(links, ",")
	w.Header().Add(Link, linkHeader)
}
