package route

import (
	"context"
	"encoding/json"
	"github.com/gobff/gobff/pkg/resource"
	"github.com/gobff/gobff/tool/cache"
	"github.com/gobff/gobff/tool/pipe"
	"github.com/gobff/gobff/tool/pipeline"
	"github.com/gobff/gobff/tool/transformer"
)

type (
	ResourceOptions struct {
		Cache       cache.Cache[json.RawMessage]
		Transformer transformer.Transformer
	}
	Resource struct {
		resource resource.Resource
		pipeline pipeline.Pipeline[json.RawMessage, json.RawMessage]
		As       string
	}
)

func NewResource(resource resource.Resource, as string, opts ResourceOptions) Resource {
	r := Resource{
		resource: resource,
		As:       as,
	}
	if opts.Cache != nil {
		r.pipeline.Add(pipe.WithCache[json.RawMessage, json.RawMessage](opts.Cache))
	}
	if opts.Transformer != nil {
		r.pipeline.Add(pipe.WithTransformer(opts.Transformer))
	}
	return r
}

func (r Resource) Run(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	return r.pipeline.Run(ctx, input, r.resource.Run)
}
