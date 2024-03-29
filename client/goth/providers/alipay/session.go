package alipay

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/admpub/goth"
	"golang.org/x/oauth2"
)

// Session stores data during the auth process with QQ.
type Session struct {
	AuthURL      string
	AccessToken  string
	RefreshToken string
	OpenID       string
	Expiry       time.Time
}

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the QQ provider.
func (s Session) GetAuthURL() (string, error) {
	if len(s.AuthURL) == 0 {
		return "", errors.New(goth.NoAuthUrlErrorMessage)
	}
	return s.AuthURL, nil
}

// Authorize the session with Github and return the access token to be stored for future use.
// ?state=state&app_id=hide&source=alipay_wallet&userOutputs=auth_user&scope=auth_user&alipay_token=&auth_code=7b7022f35fff49b896d0251bc763VX39
// documentation https://opendocs.alipay.com/apis/api_9/alipay.system.oauth.token
func (s *Session) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	p := provider.(*Provider)
	urlParams, err := p.urlParams(``, url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {params.Get("auth_code")},
		"redirect_uri": {p.CallbackURL},
	}, nil, `auth_user`)
	if err != nil {
		return ``, err
	}
	options := make([]oauth2.AuthCodeOption, 0, len(urlParams))
	for k, v := range urlParams {
		options = append(options, oauth2.SetAuthURLParam(k, v[0]))
	}
	token, err := p.config.Exchange(goth.ContextForClient(p.Client()), urlParams.Get(`code`), options...)
	if err != nil {
		return "", err
	}
	if !token.Valid() {
		return "", errors.New("Invalid token received from provider")
	}
	s.AccessToken = token.AccessToken
	s.RefreshToken = token.RefreshToken
	s.Expiry = token.Expiry
	resp, ok := token.Extra(`alipay_system_oauth_token_response`).(map[string]interface{})
	if ok {
		s.OpenID, _ = resp[`user_id`].(string)
	}
	return s.AccessToken, nil
}

// Marshal the session into a string
func (s Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s Session) String() string {
	return s.Marshal()
}

// UnmarshalSession will unmarshal a JSON string into a session.
func (p *Provider) UnmarshalSession(data string) (goth.Session, error) {
	sess := &Session{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(sess)
	return sess, err
}
