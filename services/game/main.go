package main

import (
	"flag"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"pitanti/app"
	"pitanti/conf"
	"pitanti/services/game/comp"
	"pitanti/services/game/rpc"
	"strings"
)

func main() {
	serverIdPtr := flag.Int("serverId", 1, "server id")
	envPtr := flag.String("env", "dev", "the server type")
	rpcServerPortPtr := flag.Int("rpcsvport", 0, "the port that grpc server will listen")
	flag.Parse()
	serverId, env, rpcServerPort := *serverIdPtr, *envPtr, *rpcServerPortPtr

	if rpcServerPort == 0 {
		// 通过serverId算出rpcServerPort
		rpcServerPort = serverId + 3200
	}
	appEntity := app.Create(conf.Env(env), string(conf.ServerGame), rpcServerPort, 0, serverId)
	defer appEntity.Shutdown()
	appEntity.Register(comp.NewGameComp(appEntity),
		component.WithName("comp"),
		component.WithNameFunc(strings.ToLower),
	)
	appEntity.RegisterRemote(rpc.NewGameRpc(appEntity),
		component.WithName("rpc"),
		component.WithNameFunc(strings.ToLower),
	)
	appEntity.Start()
}
