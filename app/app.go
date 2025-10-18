package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/acceptor"
	"github.com/topfreegames/pitaya/v3/pkg/cluster"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"github.com/topfreegames/pitaya/v3/pkg/config"
	"github.com/topfreegames/pitaya/v3/pkg/constants"
	"github.com/topfreegames/pitaya/v3/pkg/groups"
	"github.com/topfreegames/pitaya/v3/pkg/modules"
	router2 "pitanti/common/router"
	"pitanti/conf"
	"pitanti/services/gateway/comp"
	"pitanti/utils"
	"strconv"
	"strings"
)

func Start(env conf.Env, frontendPort int, svType string, rpcServerPort int) {
	appEntity := create(env, svType, rpcServerPort, frontendPort, 0)
	defer appEntity.Shutdown()
	if isFrontend(frontendPort) {
		configureFrontend(appEntity)
	}
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
	router := router2.NewGameRouter(appEntity)
	err := appEntity.AddRoute(string(conf.ServerGame), router.Router)
	utils.Must(err)
	if err != nil {
		fmt.Printf("error adding route %s\n", err.Error())
	}
}

func isFrontend(frontendPort int) bool {
	return frontendPort > 0
}

func create(env conf.Env, svType string, rpcServerPort, frontendPort, serverId int) (app pitaya.Pitaya) {
	if env == conf.EnvDev {
		logrus.SetLevel(logrus.DebugLevel)
	}
	isFrontend := isFrontend(frontendPort)
	meta := map[string]string{
		constants.GRPCHostKey: "127.0.0.1",
		constants.GRPCPortKey: strconv.Itoa(rpcServerPort),
		conf.CtxServerId:      strconv.Itoa(serverId),
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
