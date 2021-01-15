package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
	"github.com/sirupsen/logrus"
)

func addRequestID(config *conf.Configuration) middlewareHandler {
	return func(w http.ResponseWriter, r *http.Request) (context.Context, error) {
		id := ""
		if config.RequestID != "" {
			id = r.Header.Get(config.RequestID)
		}
		if id == "" {
			id = uuid.New().String()
		}

		ctx := r.Context()
		ctx = withRequestID(ctx, id)
		return ctx, nil
	}
}

func sendJSON(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(obj)
	if err != nil {
		err = fmt.Errorf("encoding json response %w %v", err, obj)
		return err
	}
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}

func getUserFromClaims(ctx context.Context, conn *storage.Connection) (*models.User, error) {
	claims := getClaims(ctx)
	if claims == nil {
		return nil, errors.New("Invalid token")
	}

	if claims.Subject == "" {
		return nil, errors.New("invalid claim: id")
	}

	// System Username
	if len(claims.Audience) <= 0 {
		return nil, errors.New("invalid audience")
	}

	if claims.Subject == models.SystemUserUUID.String() || claims.Subject == models.SystemUserID {
		return models.NewSystemUser(), nil
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, errors.New("Invalid user ID")
	}
	return models.FindUserByID(conn, userID)
}

func (a *API) isAdmin(ctx context.Context, u *models.User) bool {
	config := a.getConfig(ctx)
	return u.IsSuperAdmin || u.HasRole(config.JWT.AdminGroup)
}

func (a *API) getReferrer(r *http.Request) string {
	ctx := r.Context()
	config := a.getConfig(ctx)
	referrer := ""
	if reqref := r.Referer(); reqref != "" {
		base, berr := url.Parse(config.SiteURL)
		refurl, rerr := url.Parse(reqref)
		// As long as the referrer came from the site, we will redirect back there
		if berr == nil && rerr == nil && base.Hostname() == refurl.Hostname() {
			referrer = reqref
		}
	}
	return referrer
}

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"100.64.0.0/10",  // RFC6598
		"172.16.0.0/12",  // RFC1918
		"192.0.0.0/24",   // RFC6890
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

func isPrivateIP(ip net.IP) bool {
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func removeLocalhostFromPrivateIPBlock() *net.IPNet {
	_, localhost, _ := net.ParseCIDR("127.0.0.0/8")

	var localhostIndex int
	for i := 0; i < len(privateIPBlocks); i++ {
		if privateIPBlocks[i] == localhost {
			localhostIndex = i
		}
	}
	privateIPBlocks = append(privateIPBlocks[:localhostIndex], privateIPBlocks[localhostIndex+1:]...)

	return localhost
}

func unshiftPrivateIPBlock(address *net.IPNet) {
	privateIPBlocks = append([]*net.IPNet{address}, privateIPBlocks...)
}

type noLocalTransport struct {
	inner  http.RoundTripper
	errlog logrus.FieldLogger
}

func (no noLocalTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx, cancel := context.WithCancel(req.Context())

	ctx = httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{
		ConnectStart: func(network, addr string) {
			fmt.Printf("Checking network %v\n", addr)
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				cancel()
				fmt.Printf("Canceleing dur to error in addr parsing %v", err)
				return
			}
			ip := net.ParseIP(host)
			if ip == nil {
				cancel()
				fmt.Printf("Canceleing dur to error in ip parsing %v", host)
				return
			}

			if isPrivateIP(ip) {
				cancel()
				fmt.Println("Canceleing dur to private ip range")
				return
			}

		},
	})

	req = req.WithContext(ctx)
	return no.inner.RoundTrip(req)
}

func SafeRoundTripper(trans http.RoundTripper, log logrus.FieldLogger) http.RoundTripper {
	if trans == nil {
		trans = http.DefaultTransport
	}

	ret := &noLocalTransport{
		inner:  trans,
		errlog: log.WithField("transport", "local_blocker"),
	}

	return ret
}

func dataString(data map[string]interface{}, name string) string {
	val, ok := data[name]
	if !ok {
		return ""
	}
	return val.(string)
}
