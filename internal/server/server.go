package server

import (
	"encoding/json"
	"github.com/carlosrodriguesf/gobff/internal/config"
	"github.com/carlosrodriguesf/gobff/pkg/resource"
	"github.com/carlosrodriguesf/gobff/pkg/route"
	"github.com/carlosrodriguesf/gobff/pkg/source"
	"github.com/carlosrodriguesf/gobff/tool/cache"
	"github.com/carlosrodriguesf/gobff/tool/logger"
	"github.com/carlosrodriguesf/gobff/tool/transformer"
	"github.com/gin-gonic/gin"
	"strings"
)

type (
	Server interface {
		LoadConfigFile()
		Run()
	}
	serverImpl struct {
		logger      logger.Logger
		sourceMap   map[string]source.Source
		resourceMap map[string]resource.Resource
		config      *config.File
		gin         *gin.Engine
	}
)

func New() Server {
	l := logger.New()

	gin.DisableConsoleColor()
	gin.DefaultWriter = logger.NewWriterAdapter(l, logger.LevelInfo)

	return &serverImpl{
		logger:      l,
		sourceMap:   make(map[string]source.Source),
		resourceMap: make(map[string]resource.Resource),
		gin:         gin.New(),
	}
}

func (s *serverImpl) LoadConfigFile() {
	var err error
	s.config, err = config.Load(s.logger)
	if err != nil {
		s.logger.FatalE(err)
	}
}

func (s *serverImpl) Run() {
	s.instanceSources()
	s.instanceResources()
	s.instanceRoutes()
	for {
		if err := s.gin.Run("localhost:8080"); err != nil {
			s.logger.WithStackTrace().FatalE(err)
		}
		if err := recover(); err != nil {
			s.logger.ErrorE(err.(error))
		}
	}
}

func (s *serverImpl) instanceSources() {
	opts := source.Options{Logger: s.logger}
	for name, cfg := range s.config.Sources {
		src, err := source.GetSource(name, cfg.Kind, cfg.Config, opts)
		if err != nil {
			s.logger.FatalE(err)
		}
		s.sourceMap[name] = src
	}
}

func (s *serverImpl) instanceResources() {
	for name, cfg := range s.config.Resources {
		src, found := s.sourceMap[cfg.Source]
		if !found {
			s.logger.WithStackTrace().FatalF("source not found: %s", cfg.Source)
		}
		if err := src.ValidateParams(cfg.Params); err != nil {
			s.logger.WithStackTrace().FatalE(err)
		}

		var opts resource.Options
		if cfg.CacheDuration != 0 {
			opts.Cache = cache.NewCache[json.RawMessage](cfg.CacheDuration)
		}
		s.resourceMap[name] = resource.New(name, src, cfg.Params, opts)
	}
}

func (s *serverImpl) instanceRoutes() {
	for _, routeConfig := range s.config.Routes {
		routeResources := make([]route.Resource, 0)
		for resourceName, resourceConfig := range routeConfig.Resources {
			r, found := s.resourceMap[resourceName]
			if !found {
				s.logger.WithStackTrace().FatalF("resource not found: %s", resourceName)
			}
			if resourceConfig.As == "" {
				resourceConfig.As = resourceName
			}

			var opts route.ResourceOptions
			if resourceConfig.DependsOn != nil {
				opts.DependencyKeys = resourceConfig.DependsOn
			}
			if resourceConfig.Output != "" {
				t, err := transformer.New(resourceConfig.Output)
				if err != nil {
					s.logger.WithStackTrace().FatalE(err)
				}
				opts.Transformer = t
			}
			routeResources = append(routeResources, route.NewResource(r, resourceConfig.As, opts))
		}

		r := route.New(s.logger, routeResources)
		s.gin.Handle(strings.ToUpper(routeConfig.Method), routeConfig.Path, func(context *gin.Context) {
			r.ServeHTTP(context.Writer, context.Request)
		})
	}
}
