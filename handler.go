package oauth2s

import (
	"context"
	"net/http"

	"github.com/admpub/oauth2/v4/errors"
	"github.com/admpub/oauth2/v4/server"
	"github.com/webx-top/echo"
)

type HandlerInfo struct {
	PasswordAuthorization server.PasswordAuthorizationHandler
	UserAuthorize         server.UserAuthorizationHandler
	InternalError         server.InternalErrorHandler
	ResponseError         server.ResponseErrorHandler
	RefreshingScope       server.RefreshingScopeHandler
	RefreshingValidation  server.RefreshingValidationHandler
}

var (
	RequestFormDataCacheKey = `oauth2RequestForm`

	passwordAuthorizationHandler server.PasswordAuthorizationHandler = func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		err = echo.ErrNotImplemented
		return
	}

	userAuthorizeHandler server.UserAuthorizationHandler = func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		err = echo.ErrNotImplemented
		return
	}

	InternalErrorHandler server.InternalErrorHandler = func(err error) (re *errors.Response) {
		return
	}

	ResponseErrorHandler server.ResponseErrorHandler = func(re *errors.Response) {
	}
)
