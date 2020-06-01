package oauth2s

import (
	"gopkg.in/oauth2.v4/errors"
"gopkg.in/oauth2.v4/manage"
"gopkg.in/oauth2.v4/models"
"gopkg.in/oauth2.v4/server"
"gopkg.in/oauth2.v4/store"
)

func NewServer(manager *manage.Manager) (*server.Server, error) {
    // config oauth2 server
    srv := server.NewServer(server.NewConfig(), manager)
    srv.SetPasswordAuthorizationHandler(PasswordAuthorizationHandler)
    srv.SetUserAuthorizationHandler(UserAuthorizeHandler)
    srv.SetInternalErrorHandler(InternalErrorHandler)
    srv.SetResponseErrorHandler(ResponseErrorHandler)
    return srv, nil
}
