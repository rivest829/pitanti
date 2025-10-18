package main

import (
	"flag"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"pitanti/app"
	"pitanti/conf"
	gameComp "pitanti/services/game/comp"
	"strings"
)

func main() {
	svType := flag.String("type", conf.ServerGame, "the server type")
	rpcServerPort := flag.Int("rpcsvport", 3201, "the port that grpc server will listen")

	flag.Parse()
	start(svType, rpcServerPort)
}

func start(svType *string, rpcServerPort *int) {
	appEntity := app.Create(*svType, *rpcServerPort, 0)
	defer appEntity.Shutdown()

	comp := gameComp.NewGameComp(appEntity)
	appEntity.Register(comp,
		component.WithName("comp"),
		component.WithNameFunc(strings.ToLower),
	)

	appEntity.RegisterRemote(comp,
		component.WithName("rpc"),
		component.WithNameFunc(strings.ToLower),
	)
	appEntity.Start()
}
