package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wengchaoxi/auth-proxy/pkg/token"
)

const AUTH_PROXY_ENDPOINT = "/__auth_proxy__/"

type ServiceOptions struct {
	TargetURL      string
	AccessKey      string
	AuthExpiration time.Duration

	TokenManager *token.TokenManager
}

type Service struct {
	opts *ServiceOptions
}

func New(opts *ServiceOptions) *Service {
	return &Service{
		opts: opts,
	}
}

func (s *Service) Init(engine *gin.Engine) {
	auth := engine.Group(AUTH_PROXY_ENDPOINT, s.authInterrupter())
	{
		auth.Static("/", "./web")
		auth.POST("/", s.authHandler)
	}
	engine.NoRoute(s.authInterrupter(), s.proxyHandler)
}
