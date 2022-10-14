package pipeline

import (
	"context"
	"encoding/json"
	"github.com/gobff/gobff/cache"
	"github.com/gobff/gobff/transformer"
)

type (
	Next     func(ctx context.Context, input json.RawMessage) (output json.RawMessage, err error)
	Pipe     func(ctx context.Context, input json.RawMessage, next Next) (output json.RawMessage, err error)
	Pipeline struct {
		pipes []Pipe
	}
)

func (p *Pipeline) Add(pipe Pipe) {
	p.pipes = append(p.pipes, pipe)
}

func (p *Pipeline) Run(ctx context.Context, input json.RawMessage, last Next) (json.RawMessage, error) {
	return p.runIndex(ctx, input, 0, last)
}

func (p *Pipeline) runIndex(ctx context.Context, input json.RawMessage, index int, last Next) (json.RawMessage, error) {
	if index == len(p.pipes) {
		return last(ctx, input)
	}

	next := func(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
		return p.runIndex(ctx, input, index+1, last)
	}
	return p.pipes[index](ctx, input, next)
}

func WithCache(cache cache.Cache[json.RawMessage]) Pipe {
	return func(ctx context.Context, input json.RawMessage, next Next) (json.RawMessage, error) {
		if output, found := cache.Get(); found {
			return output, nil
		}

		output, err := next(ctx, input)
		if err != nil {
			return nil, err
		}

		cache.Set(output)

		return output, nil
	}
}

func WithTransformer(transformer transformer.Transformer) Pipe {
	return func(ctx context.Context, input json.RawMessage, next Next) (json.RawMessage, error) {
		output, err := next(ctx, input)
		if err != nil {
			return nil, err
		}

		output, err = transformer.Transform(output)
		if err != nil {
			return nil, err
		}

		return output, nil
	}
}
