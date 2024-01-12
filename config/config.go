package config

import (
	"flag"
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Proxy Proxy   `yaml:"proxy"`
	Auth  Auth    `yaml:"auth"`
	Rules []Rules `yaml:"rules"`
}

type Proxy struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Target string `yaml:"target"`
}

type Auth struct {
	Token string `yaml:"token"`
	OIDC  OIDC   `yaml:"oidc"`
}

type OIDC struct {
	Enabled bool `yaml:"enabled"`
}

type Rules struct {
	WhiteList []string `yaml:"whitelist"`
	BlackList []string `yaml:"blacklist"`
}

type MyConfig struct {
	config *Config
	mtx    sync.RWMutex
}

func (c *MyConfig) parse(v *viper.Viper) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	err := v.Unmarshal(c.config)
	if err != nil {
		return err
	}
	return nil
}

func (c *MyConfig) Info() Config {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return *c.config
}

func (c *MyConfig) Init(config string) error {
	v := viper.New()
	v.AutomaticEnv()
	// v.SetEnvPrefix("proxy")
	v.SetConfigFile(config)
	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	if err := c.parse(v); err != nil {
		return err
	}
	return nil
}

func New() *MyConfig {
	return &MyConfig{
		config: &Config{},
	}
}

func ParseConfig() Config {
	// config
	configFile := flag.String("config", "./config/config.yaml", "config file")
	flag.Parse()

	mc := New()
	err := mc.Init(*configFile)
	if err != nil {
		log.Panicf("init config error: %s \n", err)
	}
	return mc.Info()
}
