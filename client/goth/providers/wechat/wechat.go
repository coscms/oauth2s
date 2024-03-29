// Package wechat implements the OAuth2 protocol for authenticating users through Wechat.
// This package can be used as a reference implementation of an OAuth2 provider for Goth.
package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/admpub/goth"
	oauth2c "github.com/coscms/oauth2s/client/goth/oauth2"
	"golang.org/x/oauth2"
)

// These vars define the Authentication, Token, and API URLS for GitHub. If
// using GitHub enterprise you should change these values before calling New.
var (
	AuthURL         = "https://open.weixin.qq.com/connect/qrconnect"
	AuthURLInWechat = "https://open.weixin.qq.com/connect/oauth2/authorize"
	TokenURL        = "https://api.weixin.qq.com/sns/oauth2/access_token"
	ProfileURL      = "https://api.weixin.qq.com/sns/userinfo"
)

// New creates a new Github provider, and sets up important connection details.
// You should always call `github.New` to get a new Provider. Never try to create
// one manually.
func New(clientKey, secret, callbackURL string, scopes ...string) *Provider {
	return NewCustomisedURL(clientKey, secret, callbackURL, AuthURL, TokenURL, ProfileURL, scopes...)
}

// NewCustomisedURL is similar to New(...) but can be used to set custom URLs to connect to
func NewCustomisedURL(clientKey, secret, callbackURL, authURL, tokenURL, profileURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:    clientKey,
		Secret:       secret,
		CallbackURL:  callbackURL,
		HTTPClient:   oauth2c.DefaultClient,
		providerName: "wechat",
		profileURL:   profileURL,
	}
	p.config = newConfig(p, authURL, tokenURL, scopes)
	p.configInWechat = newConfig(p, AuthURLInWechat, tokenURL, scopes)
	return p
}

// Provider is the implementation of `goth.Provider` for accessing Github.
type Provider struct {
	ClientKey      string
	Secret         string
	CallbackURL    string
	HTTPClient     *http.Client
	config         *oauth2.Config
	configInWechat *oauth2.Config
	providerName   string
	profileURL     string
}

// Name is the name used to retrieve this provider later.
func (p *Provider) Name() string {
	return p.providerName
}

// SetName is to update the name of the provider (needed in case of multiple providers of 1 type)
func (p *Provider) SetName(name string) {
	p.providerName = name
}

func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}

func (p *Provider) urlParams(sess *Session) string {
	return `?access_token=` + url.QueryEscape(sess.AccessToken) + `&openid=` + url.QueryEscape(sess.OpenID) + `&lang=zh_CN`
}

// Debug is a no-op for the github package.
func (p *Provider) Debug(debug bool) {}

// BeginAuth asks Github for an authentication end-point.
func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	optSetter := oauth2.SetAuthURLParam(`appid`, p.ClientKey)
	session := &Session{
		AuthURL:         p.config.AuthCodeURL(state, optSetter),
		AuthURLInWechat: p.configInWechat.AuthCodeURL(state, optSetter),
	}
	return session, nil
}

// FetchUser will go to Github and access basic information about the user.
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	sess := session.(*Session)
	user := goth.User{
		AccessToken:  sess.AccessToken,
		RefreshToken: sess.RefreshToken,
		ExpiresAt:    sess.Expiry,
		Provider:     p.Name(),
	}

	if user.AccessToken == "" {
		// data is not yet retrieved since accessToken is still empty
		return user, fmt.Errorf("%s cannot get user information without accessToken", p.providerName)
	}
	if err := getOpenID(p, sess); err != nil {
		return user, err
	}

	user.UserID = sess.OpenID

	response, err := p.Client().Get(p.profileURL + p.urlParams(sess))
	if err != nil {
		return user, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return user, fmt.Errorf("Wechat API responded with a %d trying to fetch user information", response.StatusCode)
	}

	bits, err := io.ReadAll(response.Body)
	if err != nil {
		return user, err
	}

	err = json.NewDecoder(bytes.NewReader(bits)).Decode(&user.RawData)
	if err != nil {
		return user, err
	}

	err = userFromReader(bytes.NewReader(bits), &user)
	if err != nil {
		return user, err
	}

	user.RawData[`unionid`] = sess.UnionID
	return user, err
}

func userFromReader(reader io.Reader, user *goth.User) error {
	u := struct {
		Name      string `json:"nickname"`
		AvatarURL string `json:"headimgurl"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}

	user.Name = u.Name
	user.NickName = u.Name
	//user.Email = u.Email
	user.AvatarURL = u.AvatarURL
	//user.UserID = strconv.Itoa(u.ID)
	//user.Location = u.Location

	return err
}

func getOpenID(p *Provider, sess *Session) error {
	if len(sess.OpenID) > 0 {
		return nil
	}
	return errors.New(`Cannot get openid`)
}

func newConfig(provider *Provider, authURL, tokenURL string, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     provider.ClientKey,
		ClientSecret: provider.Secret,
		RedirectURL:  provider.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{},
	}

	c.Scopes = append(c.Scopes, scopes...)

	return c
}

// RefreshToken refresh token is not provided by QQ
func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	return nil, errors.New("Refresh token is not provided by Wechat")
}

// RefreshTokenAvailable refresh token is not provided by QQ
func (p *Provider) RefreshTokenAvailable() bool {
	return false
}
