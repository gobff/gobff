package pipe

import (
	"context"
	"encoding/json"
	"github.com/gobff/gobff/tool/cache"
	"github.com/gobff/gobff/tool/pipeline"
	"github.com/gobff/gobff/tool/transformer"
)

func WithCache[I, O any](cache cache.Cache[O]) pipeline.Pipe[I, O] {
	return func(ctx context.Context, input I, next pipeline.Next[I, O]) (O, error) {
		if output, found := cache.Get(); found {
			return output, nil
		}

		output, err := next(ctx, input)
		if err != nil {
			return *new(O), err
		}

		cache.Set(output)

		return output, nil
	}
}

func WithTransformer(transformer transformer.Transformer) pipeline.Pipe[json.RawMessage, json.RawMessage] {
	return func(ctx context.Context, input json.RawMessage, next pipeline.Next[json.RawMessage, json.RawMessage]) (json.RawMessage, error) {
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
