package comp

import (
	"context"
	"fmt"
	"pitanti/common/session"
	"pitanti/conf"
	"pitanti/pb"

	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/component"
)

// ConnectorRemote is a remote that will receive rpc's
type ConnectorRemote struct {
	component.Base
}

// Connector struct
type Connector struct {
	component.Base
	app pitaya.Pitaya
}

// NewConnector ctor
func NewConnector(app pitaya.Pitaya) *Connector {
	return &Connector{app: app}
}

func reply(code int32, msg string) (*pb.Response, error) {
	res := &pb.Response{
		Code: code,
		Msg:  msg,
	}
	return res, nil
}

// GetSessionData gets the session data
func (c *Connector) GetSessionData(ctx context.Context) (*session.SessionData, error) {
	s := c.app.GetSessionFromCtx(ctx)
	res := &session.SessionData{
		Data: s.GetData(),
	}
	return res, nil
}

// SetSessionData sets the session data
func (c *Connector) SetSessionData(ctx context.Context, data *session.SessionData) (*pb.Response, error) {
	s := c.app.GetSessionFromCtx(ctx)
	err := s.SetData(data.Data)
	if err != nil {
		return nil, pitaya.Error(err, "CN-000", map[string]string{"failed": "set data"})
	}
	return reply(200, "success")
}

// NotifySessionData sets the session data
func (c *Connector) NotifySessionData(ctx context.Context, data *session.SessionData) {
	s := c.app.GetSessionFromCtx(ctx)
	err := s.SetData(data.Data)
	if err != nil {
		fmt.Println("got error on notify", err)
	}
}

// SendPushToUser sends a push to a user
func (c *Connector) SendPushToUser(ctx context.Context, msg *pb.UserMessage) (*pb.Response, error) {
	_, err := c.app.SendPushToUsers("onMessage", msg, []string{"2"}, string(conf.ServerGateway))
	if err != nil {
		return nil, err
	}
	return &pb.Response{
		Code: 200,
		Msg:  "boa",
	}, nil
}

// RemoteFunc is a function that will be called remotely
func (c *ConnectorRemote) RemoteFunc(ctx context.Context, msg *pb.RPCMsg) (*pb.RPCRes, error) {
	fmt.Printf("received a remote call with this message: %s\n", msg)
	return &pb.RPCRes{
		Msg: fmt.Sprintf("received msg: %s", msg.GetMsg()),
	}, nil
}
