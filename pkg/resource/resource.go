package resource

import (
	"context"
	"encoding/json"
	"github.com/gobff/gobff/pkg/source"
	"github.com/gobff/gobff/tool/cache"
	"github.com/gobff/gobff/tool/pipe"
	"github.com/gobff/gobff/tool/pipeline"
	"github.com/gobff/gobff/tool/transformer"
	"sync"
)

type (
	Options struct {
		Cache       cache.Cache[json.RawMessage]
		Transformer transformer.Transformer
	}
	Resource interface {
		Run(ctx context.Context, input json.RawMessage) (json.RawMessage, error)
		Name() string
	}
	resource struct {
		name         string
		source       source.Source
		sourceParams source.Params
		mutex        *sync.Mutex
		pipeline     pipeline.Pipeline[json.RawMessage, json.RawMessage]
	}
)

func NewResource(name string, source source.Source, params source.Params, opts Options) Resource {
	r := &resource{
		name:         name,
		source:       source,
		sourceParams: params,
		mutex:        new(sync.Mutex),
	}
	if opts.Cache != nil {
		r.pipeline.Add(pipe.WithCache[json.RawMessage, json.RawMessage](opts.Cache))
	}
	return r
}

func (r resource) Name() string {
	return r.name
}

func (r resource) Run(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	return r.pipeline.Run(ctx, input, r.run)
}

func (r resource) run(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.source.Run(ctx, r.sourceParams, input)
}
