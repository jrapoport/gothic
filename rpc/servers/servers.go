package servers

import (
	"fmt"

	"github.com/jrapoport/gothic/api"
	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/rpc/servers/rpc"
	"github.com/jrapoport/gothic/rpc/servers/web"
	"github.com/sirupsen/logrus"
)

func ListenAndServeRPC(a *api.API, config *conf.Configuration) {
	go func() {
		addr := fmt.Sprintf("%v:%v", config.Host, config.RpcPort)
		logrus.Infof("Gothic RPC API started on: %s", addr)
		svr := rpc.NewRpcServer(a, addr)
		svr.ListenAndServe()
	}()

	go func() {
		addr := fmt.Sprintf("%v:%v", config.Host, config.RpcWebPort)
		logrus.Infof("Gothic RPC Web API started on: %s", addr)
		svr := web.NewRpcWebServer(a, addr)
		// TODO: add JWT server options
		svr.ListenAndServe()
	}()
}
