package oauth2s

import (
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/store"
)

func NewManager(config *Config) (*manage.Manager, error) {
	var err error
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	if config.Store == nil {
		config.Store, err = store.NewMemoryTokenStore()
		if err != nil {
			return nil, err
		}
	}
	manager.MapTokenStorage(config.Store)

	// access token generate method: jwt
	if config.JWTMethod != nil {
		jwtAccessGenerate := generates.NewJWTAccessGenerate(config.JWTKeyID, config.JWTKey, config.JWTMethod)
		manager.MapAccessGenerate(jwtAccessGenerate)
	}
	manager.MapClientStorage(config.ClientStore)
	return manager, err
}
