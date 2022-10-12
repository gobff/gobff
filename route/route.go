package route

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gondalf/gondalf/resource"
	"io"
	"net/http"
)

type (
	Route interface {
		Run(c *gin.Context)
	}
	Resource struct {
		resource.Resource
		Async bool
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

	response := make(map[string]Response)
	for name, res := range r.resources {
		output, err := res.Run(c, input)
		response[name] = Response{
			Data: output,
			Err:  err,
		}
	}
	c.IndentedJSON(http.StatusOK, response)
}
