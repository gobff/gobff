package source

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/carlosrodriguesf/gobff/pkg/dto"
	"github.com/carlosrodriguesf/gobff/tool/logger"
	"gopkg.in/yaml.v3"
)

const (
	KindHttp = "http"
)

type (
	Options struct {
		Logger logger.Logger
	}
	Params map[string]string
	Source interface {
		Run(ctx context.Context, srcParams Params, input dto.Request) (output json.RawMessage, err error)
		ValidateParams(params Params) error
	}
)

func GetSource(name, kind string, config yaml.Node, opts Options) (Source, error) {
	switch kind {
	case KindHttp:
		return newSourceHTTP(name, config, opts)
	default:
		return nil, fmt.Errorf("source not found: %s", kind)
	}
}
