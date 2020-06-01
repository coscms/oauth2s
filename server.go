package oauth2s

import (
    "gopkg.in/oauth2.v4/server"
)

func NewServer(config *Config) (*server.Server, error) {
    srv := server.NewServer(server.NewConfig(),config.Manager())
    passwordAuthorization := config.HandlerInfo.PasswordAuthorization
    if passwordAuthorization == nil {
        passwordAuthorization = PasswordAuthorizationHandler
    }
    srv.SetPasswordAuthorizationHandler(passwordAuthorization)
    userAuthorize := config.HandlerInfo.UserAuthorize
    if userAuthorize == nil {
        userAuthorize = UserAuthorizeHandler
    }
    srv.SetUserAuthorizationHandler(userAuthorize)
    internalError := config.HandlerInfo.InternalError
    if internalError == nil {
        internalError = InternalErrorHandler
    }
    srv.SetInternalErrorHandler(internalError)
    responseError := config.HandlerInfo.ResponseError
    if responseError == nil {
        responseError = ResponseErrorHandler
    }
    srv.SetResponseErrorHandler(ResponseErrorHandler)
    return srv, nil
}
