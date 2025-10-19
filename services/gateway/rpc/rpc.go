package rpc

import (
	"context"
	"fmt"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"pitanti/pb"
)

// ConnectorRemote is a remote that will receive rpc's
type ConnectorRemote struct {
	component.Base
}

// RemoteFunc is a function that will be called remotely
func (c *ConnectorRemote) RemoteFunc(ctx context.Context, msg *pb.RPCMsg) (*pb.RPCRes, error) {
	fmt.Printf("received a remote call with this message: %s\n", msg)
	return &pb.RPCRes{
		Msg: fmt.Sprintf("received msg: %s", msg.GetMsg()),
	}, nil
}
