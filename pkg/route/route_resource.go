package route

import (
	"context"
	"encoding/json"
	"github.com/carlosrodriguesf/gobff/pkg/dto"
	"github.com/carlosrodriguesf/gobff/pkg/resource"
	"github.com/carlosrodriguesf/gobff/tool/cache"
	"github.com/carlosrodriguesf/gobff/tool/keywatcher"
	"github.com/carlosrodriguesf/gobff/tool/syncmap"
	"github.com/carlosrodriguesf/gobff/tool/transformer"
)

type (
	r               = resource.Resource
	ResourceContext struct {
		context.Context
		watcher   keywatcher.Watcher
		resultSet syncmap.Map[ResourceResult]
	}
	ResourceOptions struct {
		Cache          cache.Cache[*ResourceResultData]
		Transformer    transformer.Transformer
		DependencyKeys []string
	}
	Resource struct {
		r
		alias          string
		cache          cache.Cache[*ResourceResultData]
		dependencyKeys []string
		transformer    transformer.Transformer
	}
)

func NewResource(resource resource.Resource, alias string, opts ResourceOptions) Resource {
	return Resource{
		r:              resource,
		alias:          alias,
		cache:          opts.Cache,
		dependencyKeys: opts.DependencyKeys,
		transformer:    opts.Transformer,
	}
}

func (r Resource) Run(ctx ResourceContext, params map[string][]string, input json.RawMessage) {
	defer ctx.watcher.Done(r.Name())

	data, err := r.runResource(ctx, params, input)
	ctx.resultSet.Set(r.Name(), ResourceResult{
		ResourceResultData: data,
		Error:              err,
		Alias:              r.alias,
	})
}

func (r Resource) runResource(ctx ResourceContext, params map[string][]string, input json.RawMessage) (*ResourceResultData, error) {
	if r.cache == nil {
		return r.runResourceWithoutCache(ctx, params, input)
	}

	if data, found := r.cache.Get(); found {
		return data, nil
	}

	data, err := r.runResourceWithoutCache(ctx, params, input)
	if err != nil {
		return nil, err
	}

	r.cache.Set(data)

	return data, nil
}

func (r Resource) runResourceWithoutCache(ctx ResourceContext, params map[string][]string, body json.RawMessage) (*ResourceResultData, error) {
	if r.dependencyKeys != nil {
		ctx.watcher.Wait(r.dependencyKeys)
	}

	data, err := r.r.Run(ctx, dto.Request{
		Params: params,
		Body:   body,
	})
	if err != nil {
		return nil, err
	}

	resultData := ResourceResultData{
		OriginData: data,
		OutputData: data,
	}

	if r.transformer != nil {
		resultData.OutputData, err = r.transformer.Transform(resultData.OriginData)
		if err != nil {
			return nil, err
		}
	}

	return &resultData, nil
}
