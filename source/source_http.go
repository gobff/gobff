package source

import (
	"bytes"
	"context"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"strings"
)

type sourceHttp struct {
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
}

func newSourceHTTP(config yaml.Node) (*sourceHttp, error) {
	var src sourceHttp
	if err := config.Decode(&src); err != nil {
		return nil, err
	}
	return &src, nil
}

func (h *sourceHttp) Run(ctx context.Context, input json.RawMessage) (output json.RawMessage, err error) {
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
