package trace

import (
	"context"
	"fmt"
	"github.com/adverax/log"
	"sync/atomic"
)

type typeTraceId int

var (
	traceIdKey typeTraceId = 1
)

type Engine struct {
	logger log.Logger
}

func New(logger log.Logger) *Engine {
	return &Engine{
		logger: logger,
	}
}

func (that *Engine) EnsureTrace(ctx context.Context, info string) context.Context {
	traceId := GetId(ctx)
	if traceId == "" {
		return that.NewTrace(ctx, info)
	}

	return ctx
}

func (that *Engine) NewTrace(ctx context.Context, info string) context.Context {
	ctx2, _ := that.NewTraceEx(ctx, "", info)
	return ctx2
}

func (that *Engine) NewTraceWithId(ctx context.Context, traceId string, info string) context.Context {
	ctx2, _ := that.NewTraceEx(ctx, traceId, info)
	return ctx2
}

func (that *Engine) NewTraceEx(ctx context.Context, traceId string, info string) (context.Context, string) {
	if traceId == "" {
		traceId = NewGUID()
	}

	if info != "" {
		that.logger.WithField("trace_id", traceId).Trace(ctx, "New trace: "+info)
	}

	return context.WithValue(ctx, traceIdKey, traceId), traceId
}

func GetId(ctx context.Context) string {
	traceId, _ := ctx.Value(traceIdKey).(string)
	return traceId
}

var NewGUID = func() func() string {
	var val int32
	return func() string {
		value := atomic.AddInt32(&val, 1)
		return fmt.Sprintf("%d", value)
	}
}()
