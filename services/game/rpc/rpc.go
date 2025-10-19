package rpc

import (
	"context"
	"fmt"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"pitanti/pb"
)

type GameRpc struct {
	component.Base
	app pitaya.Pitaya
}

// NewGame ctor
func NewGameRpc(app pitaya.Pitaya) *GameRpc {
	return &GameRpc{app: app}
}

func (c *GameRpc) XXXX(ctx context.Context, msg *pb.RPCMsg) (*pb.RPCRes, error) {
	fmt.Printf("received a remote call with this message: %s\n", msg)
	return &pb.RPCRes{
		Msg: fmt.Sprintf("received msg: %s", msg.GetMsg()),
	}, nil
}
