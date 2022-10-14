package resource

import (
	"context"
	"encoding/json"
	"github.com/gobff/gobff/cache"
	"github.com/gobff/gobff/pipeline"
	"github.com/gobff/gobff/source"
	"github.com/gobff/gobff/transformer"
	"sync"
)

type (
	Options struct {
		Cache       cache.Cache[json.RawMessage]
		Transformer transformer.Transformer
	}
	Resource interface {
		Run(ctx context.Context, input json.RawMessage) (json.RawMessage, error)
	}
	resource struct {
		source       source.Source
		sourceParams source.Params
		mutex        *sync.Mutex
		pipeline     pipeline.Pipeline
	}
)

func NewResource(source source.Source, params source.Params, opts Options) Resource {
	r := &resource{
		source:       source,
		sourceParams: params,
		mutex:        new(sync.Mutex),
	}
	if opts.Cache != nil {
		r.pipeline.Add(pipeline.WithCache(opts.Cache))
	}
	if opts.Transformer != nil {
		r.pipeline.Add(pipeline.WithTransformer(opts.Transformer))
	}
	return r
}

func (r resource) Run(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	return r.pipeline.Run(ctx, input, r.run)
}

func (r resource) run(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.source.Run(ctx, r.sourceParams, input)
}
