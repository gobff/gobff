package resource

import (
	"context"
	"encoding/json"
	"github.com/carlosrodriguesf/gobff/pkg/dto"
	"github.com/carlosrodriguesf/gobff/pkg/source"
	"github.com/carlosrodriguesf/gobff/tool/cache"
	"github.com/carlosrodriguesf/gobff/tool/logger"
	"sync"
)

type (
	Options struct {
		Logger logger.Logger
		Cache  cache.Cache[json.RawMessage]
	}
	Resource interface {
		Run(ctx context.Context, req dto.Request) (json.RawMessage, error)
		Name() string
	}
	resource struct {
		name         string
		logger       logger.Logger
		source       source.Source
		sourceParams source.Params
		mutex        *sync.Mutex
		cache        cache.Cache[json.RawMessage]
	}
)

func New(name string, source source.Source, params source.Params, opts Options) Resource {
	return &resource{
		name:         name,
		source:       source,
		sourceParams: params,
		mutex:        new(sync.Mutex),
		cache:        opts.Cache,
	}
}

func (r resource) Name() string {
	return r.name
}

func (r resource) Run(ctx context.Context, req dto.Request) (json.RawMessage, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.run(ctx, req)
}

func (r resource) run(ctx context.Context, req dto.Request) (json.RawMessage, error) {
	if r.cache == nil {
		return r.runSource(ctx, req)
	}

	if data, found := r.cache.Get(); found {
		return data, nil
	}

	data, err := r.runSource(ctx, req)
	if err != nil {
		r.logger.WithStackTrace().ErrorE(err)
		return nil, err
	}

	r.cache.Set(data)

	return data, nil
}

func (r resource) runSource(ctx context.Context, req dto.Request) (json.RawMessage, error) {
	data, err := r.source.Run(ctx, r.sourceParams, req)
	if err != nil {
		r.logger.WithStackTrace().ErrorE(err)
		return nil, err
	}
	return data, nil
}
