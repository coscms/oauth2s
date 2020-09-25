// Package alipay implements the OAuth2 protocol for authenticating users through Alipay.
// This package can be used as a reference implementation of an OAuth2 provider for Goth.
package alipay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/markbates/goth"

	"github.com/coscms/oauth2s/client/goth/oauth2"
	oauth2x "golang.org/x/oauth2"
)

// These vars define the Authentication, Token, and API URLS for GitHub. If
// using GitHub enterprise you should change these values before calling New.
var (
	SandBoxAPIURL = "https://openapi.alipaydev.com/gateway.do"

	AuthURL = "https://openauth.alipay.com/oauth2/publicAppAuthorize.htm"
	APIURL  = "https://openapi.alipay.com/gateway.do"
)

// New creates a new Github provider, and sets up important connection details.
// You should always call `github.New` to get a new Provider. Never try to create
// one manually.
func New(clientKey, secret, callbackURL string, isProduction bool, scopes ...string) *Provider {
	var apiURL string
	if isProduction {
		apiURL = APIURL
	} else {
		apiURL = SandBoxAPIURL
	}
	return NewCustomisedURL(clientKey, secret, callbackURL, apiURL, scopes...)
}

// NewCustomisedURL is similar to New(...) but can be used to set custom URLs to connect to
func NewCustomisedURL(clientKey, secret, callbackURL, apiURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:    clientKey,
		Secret:       secret,
		CallbackURL:  callbackURL,
		providerName: "alipay",
		profileURL:   apiURL,
	}
	if len(scopes) == 0 {
		scopes = []string{`auth_user`}
	}
	p.config = newConfig(p, apiURL, apiURL, scopes)
	return p
}

// Provider is the implementation of `goth.Provider` for accessing Github.
type Provider struct {
	ClientKey    string
	Secret       string
	CallbackURL  string
	HTTPClient   *http.Client
	config       *oauth2.Config
	providerName string
	profileURL   string
	debug        bool
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

func (p *Provider) createSign(params url.Values) string {
	if params.Get(`sign_type`) == `RSA` {
		return SignRSA(params, []byte(p.Secret))
	}
	return SignRSA2(params, []byte(p.Secret))
}

func (p *Provider) urlParams(method string, params url.Values, extra interface{}, scopes ...string) url.Values {
	if len(method) == 0 {
		method = `alipay.system.oauth.token`
	}
	params.Set("charset", "utf-8")
	params.Set("app_id", p.ClientKey)
	params.Set("method", method)
	params.Set("scope", strings.Join(scopes, `,`))
	params.Set("format", "JSON")
	if extra != nil {
		b, _ := json.Marshal(extra)
		params.Set("biz_content", string(b))
	}
	params.Set("sign_type", "RSA2")
	params.Set("timestamp", time.Now().Format(`2006-01-02 15:04:05`))
	params.Set("version", "1.0")
	params.Set("sign", p.createSign(params))
	return params
}

// Debug is sandbox mode
func (p *Provider) Debug(debug bool) {
	p.debug = debug
}

// BeginAuth asks Github for an authentication end-point.
func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	//url := p.config.AuthCodeURL(state)
	params := url.Values{
		"app_id":       {p.ClientKey},
		"scope":        {"auth_user"},
		"redirect_uri": {p.CallbackURL},
		"state":        {state},
	}
	session := &Session{
		AuthURL: p.config.Endpoint.AuthURL + `?` + params.Encode(),
	}
	return session, nil
}

// FetchUser will go to Github and access basic information about the user.
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	sess := session.(*Session)
	user := goth.User{
		AccessToken: sess.AccessToken,
		Provider:    p.Name(),
		RawData:     make(map[string]interface{}),
	}
	if user.AccessToken == "" {
		// data is not yet retrieved since accessToken is still empty
		return user, fmt.Errorf("%s cannot get user information without accessToken", p.providerName)
	}
	if err := getOpenID(p, sess); err != nil {
		return user, err
	}
	user.UserID = sess.OpenID

	param := url.Values{
		"auth_token": {sess.AuthCode},
	}
	param = p.urlParams(`alipay.user.info.share`, param, nil, `auth_user`)
	response, err := p.Client().Get(p.profileURL + `?` + param.Encode())
	if err != nil {
		return user, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return user, fmt.Errorf("QQ API responded with a %d trying to fetch user information", response.StatusCode)
	}

	bits, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, err
	}

	err = json.NewDecoder(bytes.NewReader(bits)).Decode(&user.RawData)
	if err != nil {
		return user, err
	}

	err = userFromReader(bytes.NewReader(bits), &user)
	return user, err
}

func userFromReader(reader io.Reader, user *goth.User) error {
	u := struct {
		Name      string `json:"nick_name"`
		AvatarURL string `json:"avatar"`
		Gender    string `json:"gender"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}

	user.Name = u.Name
	user.NickName = u.Name
	//user.Email = u.Email
	user.AvatarURL = u.AvatarURL
	user.RawData[`gender`] = u.Gender
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
	c := &oauth2x.Config{
		ClientID:     provider.ClientKey,
		ClientSecret: provider.Secret,
		RedirectURL:  provider.CallbackURL,
		Endpoint: oauth2x.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{},
	}

	for _, scope := range scopes {
		c.Scopes = append(c.Scopes, scope)
	}

	return oauth2.NewConfig(c)
}

//RefreshToken refresh token is not provided by QQ
func (p *Provider) RefreshToken(refreshToken string) (*oauth2x.Token, error) {
	return nil, errors.New("Refresh token is not provided by alipay")
}

//RefreshTokenAvailable refresh token is not provided by QQ
func (p *Provider) RefreshTokenAvailable() bool {
	return false
}
