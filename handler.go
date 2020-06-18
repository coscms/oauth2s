package oauth2s

import (
	"net/http"

	"github.com/webx-top/echo"
	"gopkg.in/oauth2.v4/errors"
)

type HandlerInfo struct {
	PasswordAuthorization func(username, password string) (userID string, err error)
	UserAuthorize         func(w http.ResponseWriter, r *http.Request) (userID string, err error)
	InternalError         func(error) *errors.Response
	ResponseError         func(*errors.Response)
}

var (
	RequestFormDataCacheKey = `oauth2RequestForm`

	PasswordAuthorizationHandler = func(username, password string) (userID string, err error) {
		return
	}

	UserAuthorizeHandler = func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		ctx := r.Context().(echo.Context)
		v := ctx.Session().Get(`uid`)
		if v == nil {
			ctx.Session().Set(RequestFormDataCacheKey, ctx.Forms())
			err = ctx.Redirect(`/oauth2/login`)
			return
		}
		userID = v.(string)
		return
	}

	InternalErrorHandler = func(err error) (re *errors.Response) {
		return
	}

	ResponseErrorHandler = func(re *errors.Response) {
	}
)
