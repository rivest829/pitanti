package comp

import (
	"context"
	"encoding/gob"
	"math/rand"
	"pitanti/pb"
	"strconv"
	"time"

	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"github.com/topfreegames/pitaya/v3/pkg/timer"
)

type (
	// GameComp represents a component that contains a bundle of game related handler
	// like Join/Message
	GameComp struct {
		component.Base
		timer *timer.Timer
		app   pitaya.Pitaya
		Stats *Stats
	}

	// Stats exports the game status
	Stats struct {
		outboundBytes int
		inboundBytes  int
	}
)

// NewGameComp returns a new game
func NewGameComp(app pitaya.Pitaya) *GameComp {
	return &GameComp{
		app:   app,
		Stats: &Stats{},
	}
}

// Init runs on service initialization
func (r *GameComp) Init() {
	r.app.GroupCreate(context.Background(), "game")
	// It is necessary to register all structs that will be used in RPC calls
	// This must be done both in the caller and callee servers
	gob.Register(&pb.UserMessage{})
}

// AfterInit component lifetime callback
func (r *GameComp) AfterInit() {
	r.timer = pitaya.NewTimer(time.Minute, func() {
		count, err := r.app.GroupCountMembers(context.Background(), "game")
		println("UserCount: Time=>", time.Now().String(), "Count=>", count, "Error=>", err)
		println("OutboundBytes", r.Stats.outboundBytes)
		println("InboundBytes", r.Stats.outboundBytes)
	})
}

func (r *GameComp) Login(ctx context.Context, msg []byte) (*pb.JoinResponse, error) {
	s := r.app.GetSessionFromCtx(ctx)
	err := s.Bind(ctx, strconv.Itoa(int(s.ID())))
	if err != nil {
		return nil, pitaya.Error(err, "RH-000", map[string]string{"failed": "bind"})
	}
	return &pb.JoinResponse{Result: "ok"}, nil
}

// Join game
func (r *GameComp) Join(ctx context.Context) (*pb.JoinResponse, error) {
	s := r.app.GetSessionFromCtx(ctx)
	s.Push("joinpush", rand.Intn(1000))
	return &pb.JoinResponse{Result: "success"}, nil
}
