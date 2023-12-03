package oauth2s

import (
	"github.com/admpub/oauth2/v4/server"
)

func NewServer(config *Config) (*server.Server, error) {
	srv := server.NewServer(server.NewConfig(), config.Manager())
	passwordAuthorization := config.HandlerInfo.PasswordAuthorization
	if passwordAuthorization == nil {
		passwordAuthorization = passwordAuthorizationHandler
	}
	srv.SetPasswordAuthorizationHandler(passwordAuthorization)
	userAuthorize := config.HandlerInfo.UserAuthorize
	if userAuthorize == nil {
		userAuthorize = userAuthorizeHandler
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
	srv.SetResponseErrorHandler(responseError)

	refreshingValidation := config.HandlerInfo.RefreshingValidation
	if refreshingValidation != nil {
		srv.SetRefreshingValidationHandler(refreshingValidation)
	}
	refreshingScope := config.HandlerInfo.RefreshingScope
	if refreshingScope != nil {
		srv.SetRefreshingScopeHandler(refreshingScope)
	}
	return srv, nil
}
