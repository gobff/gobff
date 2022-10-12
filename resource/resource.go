package resource

import (
	"context"
	"encoding/json"
	"github.com/gobff/gobff/cache"
	"github.com/gobff/gobff/source"
	"sync"
	"time"
)

type (
	Options struct {
		CacheDuration time.Duration
	}
	Resource interface {
		Run(ctx context.Context, input json.RawMessage) (json.RawMessage, error)
	}
	resource struct {
		source       source.Source
		sourceParams source.Params
		mutex        *sync.Mutex
		cache        cache.Cache[json.RawMessage]
	}
)

func NewResource(source source.Source, params source.Params, opts Options) Resource {
	if opts.CacheDuration == 0 {
		opts.CacheDuration = time.Second
	}
	r := &resource{
		source:       source,
		sourceParams: params,
		mutex:        new(sync.Mutex),
		cache:        cache.NewCache[json.RawMessage](opts.CacheDuration),
	}
	return r
}

func (r resource) Run(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if output, found := r.cache.Get(); found {
		return output, nil
	}

	output, err := r.source.Run(ctx, r.sourceParams, input)
	if err != nil {
		return nil, err
	}

	r.cache.Set(output)

	return output, err
}
