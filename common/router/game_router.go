package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/cluster"
	"github.com/topfreegames/pitaya/v3/pkg/route"
	"pitanti/conf"
	"pitanti/utils"
	"sync"
)

type GameRouter struct {
	app pitaya.Pitaya
}

func NewGameRouter(
	app pitaya.Pitaya,
) *GameRouter {

	r := &GameRouter{
		app: app,
	}
	return r
}

type LoginData struct {
	RoleId string `json:"roleId"`
}

var (
	serverIdMapLock = &sync.RWMutex{}
	serverIdMap     = make(map[uint32]*cluster.Server)
)

func (r GameRouter) Router(ctx context.Context, route *route.Route, payload []byte, servers map[string]*cluster.Server) (srv *cluster.Server, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("route role server failed, reason: %v", e))
			return
		}
	}()
	if route.String() == "game.comp.login" {
		var loginData LoginData
		err = json.Unmarshal(payload, &loginData)
		utils.Must(err)
	}
	sess := r.app.GetSessionFromCtx(ctx)
	hd := sess.GetHandshakeData()
	roleId := hd.User[conf.CtxRoleId].(string)
	serverId := utils.GetServerId(roleId)
	srv, err = BalanceServerId(serverId, servers)
	return
}

func initServerMap(servers map[string]*cluster.Server) {
	logrus.Info("init server map")
	for _, server := range servers {
		sId := utils.MustStrToInt(server.Metadata[conf.CtxServerId])
		serverIdMap[uint32(sId)] = server
	}
}

func BalanceServerId(serverId uint32, servers map[string]*cluster.Server) (targetServer *cluster.Server, err error) {
	if serverId <= 0 {
		err = errors.New("invalid_server_id")
		return
	}
	if len(serverIdMap) == 0 {
		serverIdMapLock.Lock()
		initServerMap(servers)
		serverIdMapLock.Unlock()
	}

	serverIdMapLock.RLock()
	defer serverIdMapLock.RUnlock()
	var ok bool
	targetServer, ok = serverIdMap[serverId]
	if !ok {
		err = errors.New("no_server_available")
		return
	}
	return
}
