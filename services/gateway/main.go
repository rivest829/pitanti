package main

import (
	"flag"
	"fmt"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"pitanti/app"
	router2 "pitanti/common/router"
	"pitanti/conf"
	"pitanti/services/gateway/comp"
	"pitanti/services/gateway/rpc"
	"pitanti/utils"
	"strings"
)

func main() {
	frontendPortPtr := flag.Int("port", 3100, "the port to listen")
	envPtr := flag.String("env", "dev", "the server type")
	rpcServerPortPtr := flag.Int("rpcsvport", 3200, "the port that grpc server will listen")
	flag.Parse()
	frontendPort, env, rpcServerPort := *frontendPortPtr, *envPtr, *rpcServerPortPtr
	appEntity := app.Create(conf.Env(env), string(conf.ServerGateway), rpcServerPort, frontendPort, 0)
	defer appEntity.Shutdown()
	configureFrontend(appEntity)
	appEntity.Register(comp.NewConnector(appEntity),
		component.WithName("comp"),
		component.WithNameFunc(strings.ToLower),
	)
	appEntity.RegisterRemote(&rpc.ConnectorRemote{},
		component.WithName("rpc"),
		component.WithNameFunc(strings.ToLower),
	)
	appEntity.Start()
}

// 配置前端节点特有的路由逻辑。
func configureFrontend(appEntity pitaya.Pitaya) {
	router := router2.NewGameRouter(appEntity)
	err := appEntity.AddRoute(string(conf.ServerGame), router.Router)
	utils.Must(err)
	if err != nil {
		fmt.Printf("error adding route %s\n", err.Error())
	}
}
