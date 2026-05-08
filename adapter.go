package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.gh.ink/cask/model"
	"go.gh.ink/timex"
	"go.gh.ink/toolbox/expr"
)

type Adapter struct {
	client *redis.Client
	ns     model.NamespaceInfo
}

func (a Adapter) Get(ctx context.Context) ([]byte, error) {
	return a.client.Get(ctx, a.ns.Key()).Bytes()
}

func (a Adapter) Set(ctx context.Context, value []byte, ttl timex.Duration) error {
	return a.client.Set(ctx, a.ns.Key(), value, ToSetPXDuration(ttl)).Err()
}

func (a Adapter) Del(ctx context.Context) (bool, error) {
	i, e := a.client.Del(ctx, a.ns.Key()).Result()
	return expr.Ternary(i > 0, true, false), e
}

func (a Adapter) Exists(ctx context.Context) (bool, error) {
	i, e := a.client.Exists(ctx, a.ns.Key()).Result()
	return expr.Ternary(i > 0, true, false), e
}

func (a Adapter) Expire(ctx context.Context, ttl timex.Duration) error {
	std, inf := ttl.ToStdDuration()

	switch {
	case inf == timex.PosInfTime:
		// Persist
		return a.client.Persist(ctx, a.ns.Key()).Err()
	case inf == timex.FiniteTime && std > 0:
		// Normal
		ms := std.Milliseconds()
		if ms < 1 {
			ms = 1
		}
		return a.client.PExpire(ctx, a.ns.Key(), time.Duration(ms)*time.Millisecond).Err()
	default:
		// Zero / Neg / NegInf -> Del
		return a.client.Del(ctx, a.ns.Key()).Err()
	}
}

func (a Adapter) TTL(ctx context.Context) (timex.Duration, error) {
	d, e := a.client.TTL(ctx, a.ns.Key()).Result()
	return FromTTLDuration(d), e
}
