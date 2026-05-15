package actor

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gutil"
)

type EventHandler func(msg any)

type Event struct {
	Msg     any
	Handler EventHandler
}

type Actor struct {
	buffer chan *Event
	Exit   bool
}

func NewActor(bufferSize int) *Actor {
	return &Actor{
		buffer: make(chan *Event, bufferSize),
		Exit:   false,
	}
}

// Send 发送消息
func (actor *Actor) Send(msg any, eventHandler EventHandler) {
	actor.buffer <- &Event{
		Msg:     msg,
		Handler: eventHandler,
	}
}

func (actor *Actor) Start() {
	go func() {
		gutil.TryCatch(gctx.New(), func(try context.Context) {
			for event := range actor.buffer {
				gutil.TryCatch(try, func(ctx context.Context) {
					event.Handler(event.Msg)
				}, func(ctx context.Context, exception error) {
					g.Log().Errorf(ctx, "actor 系统池 error")
				})
			}
		}, func(ctx context.Context, exception error) {
		})
	}()
}
