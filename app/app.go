package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/acceptor"
	"github.com/topfreegames/pitaya/v3/pkg/cluster"
	"github.com/topfreegames/pitaya/v3/pkg/config"
	"github.com/topfreegames/pitaya/v3/pkg/constants"
	"github.com/topfreegames/pitaya/v3/pkg/groups"
	"github.com/topfreegames/pitaya/v3/pkg/modules"
	"pitanti/conf"
	"pitanti/utils"
	"strconv"
)

// 创建并初始化一个Pitaya应用实例。
func Create(env conf.Env, svType string, rpcServerPort, frontendPort, serverId int) (app pitaya.Pitaya) {
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

// 初始化Pitaya构建器，并创建完整的Pitaya应用实例。
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

// 判断当前是否为前端节点。
func isFrontend(frontendPort int) bool {
	return frontendPort > 0
}
