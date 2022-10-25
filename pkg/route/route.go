package route

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/carlosrodriguesf/gobff/tool/keywatcher"
	"github.com/carlosrodriguesf/gobff/tool/logger"
	"github.com/carlosrodriguesf/gobff/tool/syncmap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type (
	Route interface {
		http.Handler
	}
	routeImpl struct {
		logger      logger.Logger
		resourceMap map[string]Resource
	}
)

func New(logger logger.Logger, resources []Resource) Route {
	logger = logger.AddPrefix("route")
	resourceMap := make(map[string]Resource)
	for _, resource := range resources {
		resource.setLogger(logger)
		resourceMap[resource.Name()] = resource
	}
	return &routeImpl{
		logger:      logger,
		resourceMap: resourceMap,
	}
}

func (r *routeImpl) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("context-type", "application/json")
	result := r.run(req)
	err := json.NewEncoder(writer).Encode(result)
	if err != nil {
		r.logger.ErrorE(err)
	}
}

func (r *routeImpl) getInputs(req *http.Request) (map[string]json.RawMessage, error) {
	defer req.Body.Close()
	if req.Method == http.MethodGet {
		return make(map[string]json.RawMessage), nil
	}
	var input map[string]json.RawMessage
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil && err != io.EOF {
		return nil, err
	}
	return input, nil
}

func (r *routeImpl) getParams(req *http.Request) (map[string]url.Values, error) {
	params := make(map[string]url.Values)
	for query, values := range req.URL.Query() {
		parts := strings.Split(query, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid query param: %s", query)
		}
		resName, param := parts[0], parts[1]
		if params[resName] == nil {
			params[resName] = make(url.Values)
		}
		for _, value := range values {
			params[resName].Add(param, value)
		}
	}
	return params, nil
}

func (r *routeImpl) run(req *http.Request) Result {
	inputs, err := r.getInputs(req)
	if err != nil {
		return Result{Error: err}
	}
	params, err := r.getParams(req)
	if err != nil {
		return Result{Error: err}
	}
	return Result{
		Data: r.executeResources(req.Context(), params, inputs),
	}
}

func (r *routeImpl) executeResources(ctx context.Context, params map[string]url.Values, inputs map[string]json.RawMessage) map[string]ResourceResponse {
	var (
		wg     sync.WaitGroup
		resCtx = ResourceContext{
			Context:   ctx,
			watcher:   keywatcher.New(),
			resultSet: syncmap.New[ResourceResult](),
		}
	)

	wg.Add(len(r.resourceMap))
	for name, res := range r.resourceMap {
		name, res := name, res
		go func() {
			res.Run(resCtx, params[name], inputs[name])
			wg.Done()
		}()
	}
	wg.Wait()

	return buildOutputFromResultSet(resCtx.resultSet)
}

func buildOutputFromResultSet(resultSet syncmap.Map[ResourceResult]) map[string]ResourceResponse {
	output := make(map[string]ResourceResponse)
	for _, result := range resultSet.Data() {
		if result.Omit {
			continue
		}
		if result.Error != nil {
			output[result.Alias] = ResourceResponse{Error: result.Error}
			continue
		}
		output[result.Alias] = ResourceResponse{Data: result.OutputData}
	}
	return output
}
