package echo

import (
	"net/http"
	"net/url"

    handlerIndex "github.com/admpub/webx/application/handler/frontend/index"
    modelOpen "github.com/admpub/webx/application/model/official/open"
    "github.com/coscms/oauth2s"
	"github.com/webx-top/echo"
)

func Route(router echo.IRouter) {
    g := router.Group(`/oauth2`)
	g.Route(`GET,POST`,`/authorize`, authorizeHandler)
	g.Route(`GET,POST`,`/login`, loginHandler)
	g.Route(`GET,POST`,`/logout`, logoutHandler)
    g.Route(`GET,POST`,`/token`, tokenHandler)
    g.Route(`GET,POST`,`/test`, testHandler)
}

// 首先进入执行
func authorizeHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context().(echo.Context)
    var form url.Values
    if v, y := ctx.Session().Get(oauth2s.RequestFormDataCacheKey).(map[string][]string); v != nil {
        clientID := ctx.Form("client_id")
        if len(clientID) == 0 {
            form = url.Values(v)
            r.Form = form
        }
    }

    ctx.Session().Delete(oauth2s.RequestFormDataCacheKey)

    if err := oauth2s.Default.Server().HandleAuthorizeRequest(w, r); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().(echo.Context)
    formData, ok := ctx.Session().Get(oauth2s.RequestFormDataCacheKey).(map[string][]string)
    if !ok || formData == nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    form := url.Values(formData)
    clientID := form.Get("client_id")
    if len(clientID) > 0 {
        openM := modelOpen.NewOpenApp(ctx)
        err := openM.GetAndVerifySecret(clientID, secret)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        goto PASSED
    }

	ctx.Request().Form().Set(`return_to`, `/oauth2/authorize`)
	err := handlerIndex.SignIn(ctx)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

PASSED:

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().(echo.Context)
	err := frontend.SignOut(ctx)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
    err := oauth2s.Default.Server().HandleTokenRequest(w, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func testHandler(w http.ResponseWriter, r *http.Request) {
    token, err := oauth2s.Default.Server().ValidationBearerToken(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    data := map[string]interface{}{
        "expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
        "client_id":  token.GetClientID(),
        "user_id":    token.GetUserID(),
    }
    e := json.NewEncoder(w)
    e.SetIndent("", "  ")
    e.Encode(data)
}
