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
	"pitanti/utils"
	"strings"
)

func main() {
	port := flag.Int("port", 3100, "the port to listen")
	svType := flag.String("type", conf.ServerGateway, "the server type")
	rpcServerPort := flag.Int("rpcsvport", 3200, "the port that grpc server will listen")
	flag.Parse()
	start(port, svType, rpcServerPort)
}

func start(port *int, svType *string, rpcServerPort *int) {
	appEntity := app.Create(*svType, *rpcServerPort, *port)
	defer appEntity.Shutdown()
	configureFrontend(appEntity)
	appEntity.Register(comp.NewConnector(appEntity),
		component.WithName("comp"),
		component.WithNameFunc(strings.ToLower),
	)
	appEntity.RegisterRemote(&comp.ConnectorRemote{},
		component.WithName("rpc"),
		component.WithNameFunc(strings.ToLower),
	)
	appEntity.Start()
}

func configureFrontend(appEntity pitaya.Pitaya) {

	err := appEntity.AddRoute(conf.ServerGame, router2.Router)
	utils.Must(err)

	if err != nil {
		fmt.Printf("error adding route %s\n", err.Error())
	}
}
