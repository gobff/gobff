package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gobff/gobff/internal/config"
	"github.com/gobff/gobff/pkg/resource"
	route2 "github.com/gobff/gobff/pkg/route"
	"github.com/gobff/gobff/pkg/source"
	"github.com/gobff/gobff/tool/cache"
	"github.com/gobff/gobff/tool/transformer"
	"log"
	"os"
)

type (
	Server interface {
		LoadConfigFile(path string) error
		MustLoadConfigFile(path string)
		Run() error
	}
	serverImpl struct {
		sources   map[string]source.Source
		resources map[string]resource.Resource
		config    *config.File
		gin       *gin.Engine
	}
)

func New() Server {
	return &serverImpl{
		sources:   make(map[string]source.Source),
		resources: make(map[string]resource.Resource),
		gin:       gin.New(),
	}
}

func (s *serverImpl) LoadConfigFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	s.config, err = config.Load(file)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *serverImpl) MustLoadConfigFile(path string) {
	if err := s.LoadConfigFile(path); err != nil {
		log.Fatalln(err)
	}
}

func (s *serverImpl) Run() error {
	if err := s.instanceSources(); err != nil {
		return err
	}
	if err := s.instanceResources(); err != nil {
		return err
	}
	if err := s.instanceRoutes(); err != nil {
		return err
	}
	return s.gin.Run("localhost:8080")
}

func (s *serverImpl) instanceSources() error {
	for name, cfg := range s.config.Sources {
		src, err := source.GetSource(cfg.Kind, cfg.Config)
		if err != nil {
			return err
		}
		s.sources[name] = src
	}
	return nil
}

func (s *serverImpl) instanceResources() error {
	for name, cfg := range s.config.Resources {
		src, found := s.sources[cfg.Source]
		if !found {
			return fmt.Errorf("source not found: %s", cfg.Source)
		}
		if err := src.ValidateParams(cfg.Params); err != nil {
			return err
		}

		var opts resource.Options
		if cfg.CacheDuration != 0 {
			opts.Cache = cache.NewCache[json.RawMessage](cfg.CacheDuration)
		}
		s.resources[name] = resource.NewResource(name, src, cfg.Params, opts)
	}
	return nil
}

func (s *serverImpl) instanceRoutes() error {
	for _, routeConfig := range s.config.Routes {
		routeResources := make(route2.Resources)
		for resourceName, resourceConfig := range routeConfig.Resources {
			r, found := s.resources[resourceName]
			if !found {
				return fmt.Errorf("resource not found: %s", resourceName)
			}
			if resourceConfig.As == "" {
				resourceConfig.As = resourceName
			}

			var opts route2.ResourceOptions
			if resourceConfig.DependsOn != nil {
				opts.DependsOn = resourceConfig.DependsOn
			}
			if resourceConfig.Output != "" {
				t, err := transformer.New(resourceConfig.Output)
				if err != nil {
					return err
				}
				opts.Transformer = t
			}
			routeResources[resourceName] = route2.NewResource(r, resourceConfig.As, opts)
		}
		s.gin.Handle(routeConfig.Method, routeConfig.Path, route2.New(routeConfig.Path, routeResources).Run)
	}
	return nil
}
