package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

	"github.com/wengchaoxi/auth-proxy/token"
)

const (
	RouterGroupAuth = "/auth"

	PROXY_AUTH_TOKEN = "proxy-auth-token"
	// PROXY_AUTH_ORIGIN_URL = "proxy-auth-origin-url"
)

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

var (
	expiresDuration = 1800 * time.Second
	tokenManager    = token.NewTokenManager()

	HOST       = getEnv("HOST", "localhost")
	PORT       = getEnv("PORT", "18000")
	TARGET_URL = getEnv("TARGET_URL", "http://localhost:8000")
	ACCESS_KEY = getEnv("ACCESS_KEY", "whoami")
)

func CookieInterrupter() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieToken, _ := c.Cookie(PROXY_AUTH_TOKEN)

		if !tokenManager.IsValidToken(cookieToken) {
			originalURL := url.QueryEscape(c.Request.URL.String())
			c.Redirect(http.StatusTemporaryRedirect, RouterGroupAuth+"?from="+originalURL)
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func AuthInterrupter() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieToken, _ := c.Cookie(PROXY_AUTH_TOKEN)
		if !tokenManager.IsValidToken(cookieToken) {
			c.Next()
		} else {
			targetUrl, _ := url.Parse(TARGET_URL)
			proxy := httputil.NewSingleHostReverseProxy(targetUrl)
			proxy.ServeHTTP(c.Writer, c.Request)
		}
	}
}

func main() {
	r := gin.New()
	// r.Use(gin.Logger())

	authRouter := r.Group(RouterGroupAuth, AuthInterrupter())
	{
		authRouter.GET("/", func(c *gin.Context) {
			c.File("./web/index.html")
		})
		authRouter.POST("/", func(c *gin.Context) {
			cookieToken, _ := c.Cookie(PROXY_AUTH_TOKEN)
			ak := c.PostForm("access_key")
			if ak == ACCESS_KEY || tokenManager.IsValidToken(cookieToken) {
				c.SetCookie(PROXY_AUTH_TOKEN, tokenManager.GenerateToken("whoami", expiresDuration), int(expiresDuration/time.Second), "/", c.Request.URL.Host, false, true)
				originalURL := c.Query("from")
				if originalURL == "" {
					originalURL = "/"
				} else {
					originalURL, _ = url.QueryUnescape(originalURL)
				}
				c.JSON(200, originalURL)
			} else {
				c.String(200, "0")
				c.SetCookie(PROXY_AUTH_TOKEN, "x", 0, "/", c.Request.URL.Host, false, true)
			}
		})
	}

	// proxy
	r.NoRoute(CookieInterrupter(), func(c *gin.Context) {
		targetUrl, _ := url.Parse(TARGET_URL)
		proxy := httputil.NewSingleHostReverseProxy(targetUrl)
		proxy.ServeHTTP(c.Writer, c.Request)
	})
	addr := HOST + ":" + PORT
	fmt.Printf("Running on http://%s, Proxy to %s \n", addr, TARGET_URL)
	r.Run(addr)
}
