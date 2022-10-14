package route

import (
	"context"
	"encoding/json"
	"github.com/gobff/gobff/cache"
	"github.com/gobff/gobff/pipeline"
	"github.com/gobff/gobff/resource"
	"github.com/gobff/gobff/transformer"
)

type (
	ResourceOptions struct {
		Cache       cache.Cache[json.RawMessage]
		Transformer transformer.Transformer
	}
	Resource struct {
		resource resource.Resource
		pipeline pipeline.Pipeline
		As       string
	}
)

func NewResource(resource resource.Resource, as string, opts ResourceOptions) Resource {
	r := Resource{
		resource: resource,
		As:       as,
	}
	if opts.Cache != nil {
		r.pipeline.Add(pipeline.WithCache(opts.Cache))
	}
	if opts.Transformer != nil {
		r.pipeline.Add(pipeline.WithTransformer(opts.Transformer))
	}
	return r
}

func (r Resource) Run(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	return r.pipeline.Run(ctx, input, r.resource.Run)
}
