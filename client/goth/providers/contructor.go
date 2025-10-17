package providers

import (
	"github.com/admpub/goth"

	"github.com/webx-top/echo/handler/oauth2"

	"github.com/coscms/oauth2s/client/goth/providers/alipay"
	"github.com/coscms/oauth2s/client/goth/providers/qq"
	"github.com/coscms/oauth2s/client/goth/providers/wechat"
	"github.com/coscms/oauth2s/client/goth/providers/weibo"
)

func init() {
	oauth2.Register(`alipay`, func(account *oauth2.Account) goth.Provider {
		return alipay.New(account.Key, account.Secret, account.CallbackURL, true)
	})
	oauth2.Register(`alipay_dev`, func(account *oauth2.Account) goth.Provider {
		return alipay.New(account.Key, account.Secret, account.CallbackURL, false)
	})
	oauth2.Register(`qq`, func(account *oauth2.Account) goth.Provider {
		return qq.New(account.Key, account.Secret, account.CallbackURL)
	})
	oauth2.Register(`weibo`, func(account *oauth2.Account) goth.Provider {
		return weibo.New(account.Key, account.Secret, account.CallbackURL)
	})
	oauth2.Register(`wechat`, func(account *oauth2.Account) goth.Provider {
		return wechat.New(account.Key, account.Secret, account.CallbackURL)
	})
}
