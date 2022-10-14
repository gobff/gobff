package route

import (
	"context"
	"encoding/json"
	"github.com/gobff/gobff/pkg/resource"
	"github.com/gobff/gobff/tool/cache"
	"github.com/gobff/gobff/tool/donewatcher"
	"github.com/gobff/gobff/tool/pipe"
	"github.com/gobff/gobff/tool/pipeline"
	"github.com/gobff/gobff/tool/syncmap"
	"github.com/gobff/gobff/tool/transformer"
)

type (
	ResourceOptions struct {
		Cache       cache.Cache[json.RawMessage]
		Transformer transformer.Transformer
		DependsOn   []string
	}
	Resource struct {
		alias     string
		dependsOn []string
		resource  resource.Resource
		pipeline  pipeline.Pipeline[json.RawMessage, json.RawMessage]
	}
)

func NewResource(resource resource.Resource, alias string, opts ResourceOptions) Resource {
	r := Resource{
		resource:  resource,
		alias:     alias,
		dependsOn: opts.DependsOn,
	}
	if opts.Cache != nil {
		r.pipeline.Add(pipe.WithCache[json.RawMessage, json.RawMessage](opts.Cache))
	}
	if opts.Transformer != nil {
		r.pipeline.Add(pipe.WithTransformer(opts.Transformer))
	}
	return r
}

func (r Resource) Run(ctx context.Context, input json.RawMessage, responseMap syncmap.Map[Response], watcher donewatcher.Watcher) {
	defer watcher.Done(r.resource.Name())

	if r.dependsOn != nil {
		watcher.Wait(r.dependsOn)
	}

	output, err := r.pipeline.Run(ctx, input, r.resource.Run)
	responseMap.Set(r.alias, Response{
		Data: output,
		Err:  err,
	})
}
