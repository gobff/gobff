package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gondalf/gondalf/resource"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
)

type Http struct {
	Name   string
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
}

func (h *Http) Run(ctx context.Context, input json.RawMessage) (output json.RawMessage, err error) {
	req, err := http.NewRequest(h.Method, h.Path, bytes.NewReader(input))
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

func FactoryFunc(name string, config yaml.Node) (resource.Resource, error) {
	res := &Http{
		Name: name,
	}
	if err := config.Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}
