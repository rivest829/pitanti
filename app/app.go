package app

import (
	"fmt"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/acceptor"
	"github.com/topfreegames/pitaya/v3/pkg/cluster"
	"github.com/topfreegames/pitaya/v3/pkg/config"
	"github.com/topfreegames/pitaya/v3/pkg/constants"
	"github.com/topfreegames/pitaya/v3/pkg/groups"
	"github.com/topfreegames/pitaya/v3/pkg/modules"
	"pitanti/utils"
	"strconv"
)

func Create(svType string, rpcServerPort int, frontendPort int) (app pitaya.Pitaya) {
	isFrontend := frontendPort > 0
	meta := map[string]string{
		constants.GRPCHostKey: "127.0.0.1",
		constants.GRPCPortKey: strconv.Itoa(rpcServerPort),
	}

	var bs *modules.ETCDBindingStorage
	app, bs = initPitaya(frontendPort, isFrontend, svType, meta, rpcServerPort)

	err := app.RegisterModule(bs, "bindingsStorage")
	utils.Must(err)

	return
}

func initPitaya(port int, isFrontend bool, svType string, meta map[string]string, rpcServerPort int) (pitaya.Pitaya, *modules.ETCDBindingStorage) {
	builder := pitaya.NewDefaultBuilder(isFrontend, svType, pitaya.Cluster, meta, *config.NewDefaultPitayaConfig())

	grpcServerConfig := builder.Config.Cluster.RPC.Server.Grpc
	grpcServerConfig.Port = rpcServerPort
	gs, err := cluster.NewGRPCServer(grpcServerConfig, builder.Server, builder.MetricsReporters)
	if err != nil {
		panic(err)
	}
	builder.RPCServer = gs
	builder.Groups = groups.NewMemoryGroupService(builder.Config.Groups.Memory)

	bs := modules.NewETCDBindingStorage(builder.Server, builder.SessionPool, builder.Config.Modules.BindingStorage.Etcd)

	gc, err := cluster.NewGRPCClient(
		builder.Config.Cluster.RPC.Client.Grpc,
		builder.Server,
		builder.MetricsReporters,
		bs,
		cluster.NewInfoRetriever(builder.Config.Cluster.Info),
	)
	if err != nil {
		panic(err)
	}
	builder.RPCClient = gc

	if isFrontend {
		tcp := acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", port))
		builder.AddAcceptor(tcp)
	}

	return builder.Build(), bs
}
