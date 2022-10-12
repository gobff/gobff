package source

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"strings"
)

type sourceHttp struct {
	URL string `yaml:"url"`
}

func newSourceHTTP(config yaml.Node) (*sourceHttp, error) {
	var src sourceHttp
	if err := config.Decode(&src); err != nil {
		return nil, err
	}
	return &src, nil
}

func (h *sourceHttp) Run(ctx context.Context, params Params, input json.RawMessage) (output json.RawMessage, err error) {
	req, err := http.NewRequest(strings.ToUpper(params["method"]), h.URL+params["path"], bytes.NewReader(input))
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

func (h *sourceHttp) ValidateParams(params Params) error {
	if params["method"] == "" {
		return errors.New("param 'method' for http source is required")
	}
	if params["path"] == "" {
		return errors.New("param 'path' for http source is required")
	}
	return nil
}
