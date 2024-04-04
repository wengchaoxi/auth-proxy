package service

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	AUTH_PROXY_TOKEN_NAME = "auth_proxy_token"
	// AUTH_PROXY_ORIGIN_URL = "auth-proxy-origin-url"
)

func (s *Service) authInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieToken, _ := c.Cookie(AUTH_PROXY_TOKEN_NAME)
		if !s.opts.TokenManager.IsValidToken(cookieToken) {
			reqPath := c.Request.URL.Path
			if reqPath == AUTH_PROXY_ENDPOINT {
				c.Next()
			} else {
				location := AUTH_PROXY_ENDPOINT
				reqUrlStr := c.Request.URL.String()
				if reqUrlStr != "/" {
					location = location + "?from=" + url.QueryEscape(reqUrlStr)
				}
				c.Redirect(http.StatusTemporaryRedirect, location)
				c.Abort()
			}
		} else {
			c.Next()
		}
	}
}

func (s *Service) authHandler(c *gin.Context) {
	cookieToken, _ := c.Cookie(AUTH_PROXY_TOKEN_NAME)
	ak := c.PostForm("access_key")
	if ak == s.opts.AccessKey || s.opts.TokenManager.IsValidToken(cookieToken) {
		expiresAt := time.Now().Add(s.opts.AuthExpiration).Unix()
		c.SetCookie(AUTH_PROXY_TOKEN_NAME,
			s.opts.TokenManager.GenerateToken(AUTH_PROXY_TOKEN_NAME, expiresAt),
			int(s.opts.AuthExpiration/time.Second), "/", c.Request.URL.Host, false, true)
		originalURL := c.Query("from")
		if originalURL == "" {
			originalURL = "/"
		} else {
			originalURL, _ = url.QueryUnescape(originalURL)
		}
		c.JSON(200, originalURL)
	} else {
		c.String(200, "0")
		c.SetCookie(AUTH_PROXY_TOKEN_NAME, "x", 0, "/", c.Request.URL.Host, false, true)
	}
}

func (s *Service) proxyHandler(c *gin.Context) {
	targetUrl, _ := url.Parse(s.opts.TargetURL)
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}
