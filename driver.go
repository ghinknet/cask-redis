package redis

import (
	"github.com/redis/go-redis/v9"
	"go.gh.ink/cask/model"
)

type Driver struct{}

func (d Driver) NewAdapter(client any, ns model.NamespaceInfo) (adapter model.Adapter, ok bool) {
	if c, ok := client.(*redis.Client); ok && c != nil {
		return Adapter{
			client: c,
			ns:     ns,
		}, true
	}
	return nil, false
}
