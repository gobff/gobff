package source

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/carlosrodriguesf/gobff/pkg/dto"
	"github.com/carlosrodriguesf/gobff/tool/logger"
	"gopkg.in/yaml.v3"
	"net/http"
	"strings"
)

type (
	sourceHttpConfig struct {
		URL     string            `yaml:"url"`
		Headers map[string]string `yaml:"headers"`
	}
	sourceHttp struct {
		logger  logger.Logger
		url     string
		headers map[string]string
	}
)

func newSourceHTTP(name string, config yaml.Node, opts Options) (*sourceHttp, error) {
	var srcConfig sourceHttpConfig
	if err := config.Decode(&srcConfig); err != nil {
		return nil, err
	}
	return &sourceHttp{
		logger:  opts.Logger.WithPrefix("source.http." + name),
		url:     srcConfig.URL,
		headers: srcConfig.Headers,
	}, nil
}

func (h *sourceHttp) Run(ctx context.Context, params Params, input dto.Request) (output json.RawMessage, err error) {
	req, err := h.buildRequest(ctx, params, input)
	if err != nil {
		h.logger.WithStackTrace().ErrorE(err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		h.logger.WithStackTrace().ErrorE(err)
		return nil, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			h.logger.WithStackTrace().ErrorE(err)
		}
	}()

	err = json.NewDecoder(res.Body).Decode(&output)
	if err != nil {
		h.logger.WithStackTrace().ErrorE(err)
		return nil, err
	}

	return output, nil
}

func (h *sourceHttp) ValidateParams(params Params) error {
	if params["method"] == "" {
		return errors.New("param 'method' for http source is required")
	}
	if params["path"] == "" {
		return errors.New("param 'path' for http source is required")
	}
	return nil
}

func (h *sourceHttp) buildRequest(ctx context.Context, params Params, input dto.Request) (*http.Request, error) {
	requestURL := h.buildURL(params, input)
	req, err := http.NewRequest(strings.ToUpper(params["method"]), requestURL, bytes.NewReader(input.Body))
	if err != nil {
		return nil, err
	}
	for header, value := range h.headers {
		req.Header.Set(header, value)
	}
	return req.WithContext(ctx), nil
}

func (h *sourceHttp) buildURL(params Params, input dto.Request) string {
	var builder strings.Builder
	builder.WriteString(h.url)
	builder.WriteString(params["path"])
	if input.Params != nil {
		builder.WriteString("?")
		builder.WriteString(input.Params.Encode())
	}
	return builder.String()
}
