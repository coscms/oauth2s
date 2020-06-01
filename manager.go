package oauth2s

import (
	"gopkg.in/oauth2.v4/errors"
	"gopkg.in/oauth2.v4/manage"
	"gopkg.in/oauth2.v4/models"
	"gopkg.in/oauth2.v4/server"
	"gopkg.in/oauth2.v4/store"
	"gopkg.in/oauth2.v4/generates"
	
    "github.com/dgrijalva/jwt-go"
    "github.com/go-redis/redis"
    oredis "gopkg.in/go-oauth2/redis.v4"
)

func NewManager(store store.Store) *Manager {
    // manager config
    manager := manage.NewDefaultManager()
    manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
    // token store
    // manager.MustTokenStorage(store.NewMemoryTokenStore())
    // use redis token store
    manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
        Addr: yaml.Cfg.Redis.Default.Addr,
        DB: yaml.Cfg.Redis.Default.Db,
    }))

    // access token generate method: jwt
    manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("00000000"), jwt.SigningMethodHS512))
    clientStore := store.NewClientStore()
    for _, v := range yaml.Cfg.OAuth2.Client {
        clientStore.Set(v.ID, &models.Client{
            ID:     v.ID,
            Secret: v.Secret,
            Domain: v.Domain,
        })
    }
    manager.MapClientStorage(clientStore)
}