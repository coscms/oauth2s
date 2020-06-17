package echoauth

import (
	"github.com/coscms/oauth2s"
	adapter "github.com/coscms/oauth2s/adapter/echo"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/oauth2.v4/store"
	"gopkg.in/oauth2.v4/models"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/engine"
)

func init() {
	clientStore := store.NewClientStore()
	clientStore.Set("222222", &models.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:9094",
	})
	oauth2s.Default.Init(
		oauth2s.JWTKey([]byte("00000000")),
		oauth2s.JWTMethod(jwt.SigningMethodHS512),
		oauth2s.ClientStore(clientStore),
	)
}

func main() {
	adapter.Route(defaults.Default)
	c := &engine.Config{
		Address:     ":9094",
		TLSAuto:     false,
	}
	eng := standard.New()
	defaults.Default.Run(standard.NewWithConfig(c))
}