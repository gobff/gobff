package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gondalf/gondalf/source"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"strings"
)

type Http struct {
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
}

func (h *Http) Run(ctx context.Context, input json.RawMessage) (output json.RawMessage, err error) {
	req, err := http.NewRequest(strings.ToUpper(h.Method), h.Path, bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("error: %v\n", err)
		}
	}()

	return output, json.NewDecoder(res.Body).Decode(&output)
}

func FactoryFunc(config yaml.Node) (source.Source, error) {
	var src Http
	if err := config.Decode(&src); err != nil {
		return nil, err
	}
	return &src, nil
}
