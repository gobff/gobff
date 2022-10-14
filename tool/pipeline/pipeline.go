package pipeline

import (
	"context"
)

type (
	Next[I, O any]     func(ctx context.Context, input I) (output O, err error)
	Pipe[I, O any]     func(ctx context.Context, input I, next Next[I, O]) (output O, err error)
	Pipeline[I, O any] struct {
		pipes []Pipe[I, O]
	}
)

func (p *Pipeline[I, O]) Add(pipe Pipe[I, O]) {
	p.pipes = append(p.pipes, pipe)
}

func (p *Pipeline[I, O]) Run(ctx context.Context, input I, last Next[I, O]) (O, error) {
	return p.runIndex(ctx, input, 0, last)
}

func (p *Pipeline[I, O]) runIndex(ctx context.Context, input I, index int, last Next[I, O]) (O, error) {
	if index == len(p.pipes) {
		return last(ctx, input)
	}

	next := func(ctx context.Context, input I) (O, error) {
		return p.runIndex(ctx, input, index+1, last)
	}
	return p.pipes[index](ctx, input, next)
}
