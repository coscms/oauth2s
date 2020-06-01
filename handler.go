package oauth2s

import (
	"net/http"

	"github.com/webx-top/echo"
	"gopkg.in/oauth2.v4/errors"
	"github.com/llaoj/oauth2/utils/log"
)

type HandlerInfo struct{
	PasswordAuthorization func(username, password string) (userID string, err error)
	UserAuthorize func(w http.ResponseWriter, r *http.Request) (userID string, err error)
	InternalError func (error) (*errors.Response)
	ResponseError func (*errors.Response)
}

var (
	PasswordAuthorizationHandler = func(username, password string) (userID string, err error) {
    	var user model.User
    	userID = user.GetUserIDByPwd(username, password)
    	return
	}

	UserAuthorizeHandler = func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
    	v, _ := session.Get(r, "LoggedInUserID")
    	if v == nil {
       		if r.Form == nil {
            	r.ParseForm()
        	}
        	err = session.Set(w, r, "RequestForm", r.Form)
        	if err != nil {
        	    log.App.Error(err.Error())
        	    return
        	}
		
        	w.Header().Set("Location", "/login")
        	w.WriteHeader(http.StatusFound)

        	return
    	}
    	userID = v.(string)

    	// 不记住用户
    	// store.Delete("LoggedInUserID")
    	// store.Save()

    	return
	}

	InternalErrorHandler = func (err error) (re *errors.Response) {
    	log.App.Error("Internal Error:", err.Error())
    
    	return
	}

	ResponseErrorHandler = func (re *errors.Response) {
    	log.App.Error("Response Error:", re.Error.Error())
	}
)

func Route(router echo.IRouter) {
	router.Route(`GET,POST`,`/authorize`, authorizeHandler)
	router.Route(`GET,POST`,`/login`, loginHandler)
	router.Route(`GET,POST`,`/logout`, logoutHandler)
	router.Route(`GET,POST`,`/token`, tokenHandler)
}
