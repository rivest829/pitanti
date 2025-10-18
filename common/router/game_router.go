package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/topfreegames/pitaya/v3/pkg/cluster"
	"github.com/topfreegames/pitaya/v3/pkg/route"
	"github.com/zehuamama/balancer/balancer"
	"strconv"
)

func Router(ctx context.Context, route *route.Route, payload []byte, servers map[string]*cluster.Server) (srv *cluster.Server, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("route role server failed, reason: %v", e))
			return
		}
	}()
	srv, err = BalanceServerId(1, servers)
	return
}

func BalanceServerId(serverId int32, servers map[string]*cluster.Server) (server *cluster.Server, err error) {
	if serverId <= 0 {
		err = errors.New("invalid_server_id")
		return
	}
	var serverIds []string
	for _, srv := range servers {
		serverIds = append(serverIds, srv.ID)
	}
	b, err := balancer.Build(balancer.ConsistentHashBalancer, serverIds)
	if err != nil {
		return
	}
	sdServerId, err := b.Balance(strconv.Itoa(int(serverId)))
	if err != nil {
		return
	}
	var ok bool
	server, ok = servers[sdServerId]
	if !ok {
		err = errors.New("no_server_available")
		return
	}
	return
}
