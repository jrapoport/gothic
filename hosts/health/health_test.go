package health

import (
	"net/http"
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
)

const testResponse = `{"name":"gothic","version":"debug","status":"bela lugos` +
	`i's dead","hosts":{"test1":{"name":"test1","online":false},"test2":{"nam` +
	`e":"test2","address":"` + tsrv.MockAddress + `","online":true},"test3":{"n` +
	`ame":"mock","address":"` + tsrv.MockAddress + `","online":true}}}`

func TestNewHealthHost(t *testing.T) {
	tests := []struct {
		h      core.Hosted
		name   string
		addr   string
		online bool
	}{
		{nil, "", "", false},
		{&tsrv.MockHost{}, "", tsrv.MockAddress, true},
		{tsrv.NewMockHost("mock"), "mock", tsrv.MockAddress, true},
	}
	for _, test := range tests {
		stat := hostStatus(test.h)
		assert.Equal(t, test.name, stat.Name)
		assert.Equal(t, test.addr, stat.Address)
		assert.Equal(t, test.online, stat.Online)
	}
	hosted := map[string]core.Hosted{
		"test1": tests[0].h,
		"test2": tests[1].h,
		"test3": tests[2].h,
	}
	a, c, _ := tcore.API(t, false)
	srv := rest.NewServer(core.NewServer(a, "test"))
	rt := rest.NewRouter(c)
	s := &http.Server{Handler: rt}
	registerServer(s, srv, hosted)
	assert.HTTPBodyContains(t, s.Handler.ServeHTTP,
		http.MethodGet, config.HealthEndpoint, nil, testResponse)
}
