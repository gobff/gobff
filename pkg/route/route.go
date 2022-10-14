package route

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gobff/gobff/tool/donewatcher"
	"github.com/gobff/gobff/tool/syncmap"
	"io"
	"net/http"
	"sync"
)

type (
	Route interface {
		Run(c *gin.Context)
	}
	Resources map[string]Resource
	routeImpl struct {
		path      string
		resources Resources
	}
)

func New(path string, resources Resources) Route {
	return &routeImpl{
		path:      path,
		resources: resources,
	}
}

func (r *routeImpl) Run(c *gin.Context) {
	var input json.RawMessage
	if err := c.Bind(&input); err != nil && err != io.EOF {
		c.JSON(http.StatusInternalServerError, Response{Err: err})
		return
	}
	c.IndentedJSON(
		http.StatusOK,
		r.executeResources(c, input),
	)
}

func (r *routeImpl) executeResources(ctx context.Context, input json.RawMessage) map[string]Response {
	var (
		wg          sync.WaitGroup
		responseMap = syncmap.New[Response]()
		watcher     = donewatcher.NewWatcher()
	)

	wg.Add(len(r.resources))
	for _, res := range r.resources {
		res := res
		go func() {
			res.Run(ctx, input, responseMap, watcher)
			wg.Done()
		}()
	}
	wg.Wait()

	return responseMap.Data()
}
