package echo

import (
	"net/http"
	"fmt"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/coscms/oauth2s"
    modelOpen "github.com/admpub/webx/application/model/official/open"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/admpub/webx/application/dbschema"
)

var loginVerifyings = map[string]func(ctx echo.Context)(userID string, err error){
	`password`:func(ctx echo.Context) (userID string, err error) {
		username := ctx.Form(`username`)
		password := ctx.Form(`password`)
		if username == "user" && password == "userpwd" {
			ctx.Session().Set(`LoggedInUserID`, `1`)
			return `1`, nil
		}
		err = ctx.E(`账号或密码错误！`)
		return
	},
}

func PasswordAuthorizationHandler(username, password string) (userID string, err error) {
	c := defaults.NewMockContext()
	m := modelCustomer.NewCustomer(c)
	var authType string
	err = m.SignIn(username, username, authType)
	if err != nil {
		return
	}
	userID = fmt.Sprint(m.Id)
	return
}

func UserAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	ctx := r.Context().(echo.Context)
	customer, ok := ctx.Session().Get(`customer`).(*dbschema.OfficialCustomer)
	if !ok || customer == nil {
		ctx.Session().Set(oauth2s.RequestFormDataCacheKey, ctx.Request().Forms())
		err = ctx.Redirect(`/oauth2/login`)
		return
	}
	userID = fmt.Sprint(customer.Id)
	//ctx.Session().Delete(`customer`)
	return
}

func InternalErrorHandler(err error) (re *errors.Response) {
	return
}

func ResponseErrorHandler(re *errors.Response) {
}

func init() {
	oauth2.PasswordAuthorizationHandler = PasswordAuthorizationHandler
	oauth2.UserAuthorizeHandler = UserAuthorizeHandler
	oauth2.InternalErrorHandler = InternalErrorHandler
	oauth2.ResponseErrorHandler = ResponseErrorHandler
}
