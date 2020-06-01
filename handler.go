package oauth2s

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