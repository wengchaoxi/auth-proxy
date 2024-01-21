package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/wengchaoxi/auth-proxy/pkg/token"
	"github.com/wengchaoxi/auth-proxy/service"
)

var (
	HOST            = getEnv("HOST", "localhost")
	PORT            = getEnv("PORT", "18000")
	TARGET_URL      = getEnv("TARGET_URL", "https://github.com/wengchaoxi/auth-proxy")
	ACCESS_KEY      = getEnv("AUTH_ACCESS_KEY", "whoami")
	AUTH_EXPIRATION = getEnv("AUTH_EXPIRATION", "24h")
	authExpiration  = func() time.Duration {
		if d, err := time.ParseDuration(AUTH_EXPIRATION); err == nil {
			return d
		}
		return 24 * time.Hour
	}()
)

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func generateSecretKey(x string) string {
	size := len(x)
	if size > 0 {
		for {
			if len(x) >= 32 {
				return x[:32]
			}
			x += x
		}
	}
	return "passphrasewhichneedstobe32bytes!"
}

func main() {
	app := service.New(&service.ServiceOptions{
		TargetURL:      TARGET_URL,
		AccessKey:      ACCESS_KEY,
		AuthExpiration: authExpiration,
		TokenManager:   token.NewTokenManager(generateSecretKey(ACCESS_KEY)),
	})
	engine := gin.Default()
	app.Init(engine)

	var err error
	var wg sync.WaitGroup
	addr := HOST + ":" + PORT
	srv := &http.Server{Addr: addr, Handler: engine}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("server err: %v\n", err)
			os.Exit(-1)
		}
	}()
	log.Printf("Running on http://%s, Proxy to %s\n", addr, TARGET_URL)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown err: %v\n", err)
	}
	wg.Wait()
}
