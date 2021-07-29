package oauth2s

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/golang-jwt/jwt"
)

func NewConfig() *Config {
	return &Config{
		JWTMethod:   jwt.SigningMethodHS512,
		HandlerInfo: &HandlerInfo{},
	}
}

type Config struct {
	JWTKeyID    string
	JWTKey      []byte
	JWTMethod   jwt.SigningMethod
	Store       oauth2.TokenStore
	ClientStore oauth2.ClientStore
	HandlerInfo *HandlerInfo
	manager     *manage.Manager
	server      *server.Server
}

func (c *Config) InitConfig(options ...OptionsSetter) *Config {
	for _, fn := range options {
		fn(c)
	}
	return c
}

func (c *Config) Init(options ...OptionsSetter) error {
	var err error
	c.InitConfig(options...)
	c.manager, err = NewManager(c)
	if err != nil {
		return err
	}
	c.server, err = NewServer(c)
	return err
}

func (c *Config) Manager() *manage.Manager {
	return c.manager
}

func (c *Config) Server() *server.Server {
	return c.server
}

type OptionsSetter func(*Config)

func JWTKeyID(keyID string) OptionsSetter {
	return func(c *Config) {
		c.JWTKeyID = keyID
	}
}

func JWTKey(key []byte) OptionsSetter {
	return func(c *Config) {
		c.JWTKey = key
	}
}

func JWTMethod(method jwt.SigningMethod) OptionsSetter {
	return func(c *Config) {
		c.JWTMethod = method
	}
}

func SetStore(store oauth2.TokenStore) OptionsSetter {
	return func(c *Config) {
		c.Store = store
	}
}

func ClientStore(store oauth2.ClientStore) OptionsSetter {
	return func(c *Config) {
		c.ClientStore = store
	}
}

func SetHandler(handlerInfo *HandlerInfo) OptionsSetter {
	return func(c *Config) {
		c.HandlerInfo = handlerInfo
	}
}
