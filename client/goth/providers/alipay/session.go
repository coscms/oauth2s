package alipay

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/coscms/oauth2s/client/goth/oauth2"
	"github.com/markbates/goth"
)

// Session stores data during the auth process with QQ.
type Session struct {
	AuthURL     string
	AuthCode    string
	AccessToken string
	OpenID      string
}

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the QQ provider.
func (s Session) GetAuthURL() (string, error) {
	if len(s.AuthURL) == 0 {
		return "", errors.New(goth.NoAuthUrlErrorMessage)
	}
	return s.AuthURL, nil
}

// Authorize the session with Github and return the access token to be stored for future use.
func (s *Session) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	p := provider.(*Provider)
	values := params.(url.Values)
	if err := VerifySign(values, []byte(p.Secret)); err != nil {
		return ``, err
	}
	s.AuthCode = params.Get("auth_code")
	urlParams := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {s.AuthCode},
		"redirect_uri": oauth2.CondVal(p.CallbackURL),
	}
	urlParams = p.urlParams(``, urlParams, nil, `auth_user`)
	token, err := p.config.Exchange(goth.ContextForClient(p.Client()), urlParams)
	if err != nil {
		return "", err
	}
	if !token.Valid() {
		return "", errors.New("Invalid token received from provider")
	}
	s.AccessToken = token.AccessToken
	if r, y := token.Raw[`alipay_system_oauth_token_response`]; y {
		if m, y := r.(map[string]interface{}); y {
			if v, y := m[`user_id`]; y {
				s.OpenID, _ = v.(string)
			}
		}
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
