package resource

import (
	"context"
	"encoding/json"
	"github.com/gondalf/gondalf/source"
	"github.com/jellydator/ttlcache/v3"
	"sync"
)

type (
	Resource interface {
		Run(ctx context.Context, input json.RawMessage) (output json.RawMessage, err error)
	}
	resource struct {
		source  source.Source
		mutex   *sync.Mutex
		cache   ttlcache.Cache[string, json.RawMessage]
		running bool
	}
)

func NewResource(source source.Source) Resource {
	r := &resource{source: source}
	return r
}

func (r resource) Run(ctx context.Context, input json.RawMessage) (output json.RawMessage, err error) {
	return r.source.Run(ctx, input)
}
