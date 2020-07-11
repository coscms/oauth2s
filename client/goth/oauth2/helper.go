package oauth2

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
)

func CondVal(v string) []string {
	if len(v) == 0 {
		return nil
	}
	return []string{v}
}

// HTTPClient is the context key to use with golang.org/x/net/context's
// WithValue function to associate an *http.Client value with a context.
var HTTPClient ContextKey

// ContextKey is just an empty struct. It exists so HTTPClient can be
// an immutable public variable with a unique type. It's immutable
// because nobody else can create a ContextKey, being unexported.
type ContextKey struct{}

// ContextClientFunc is a func which tries to return an *http.Client
// given a Context value. If it returns an error, the search stops
// with that error.  If it returns (nil, nil), the search continues
// down the list of registered funcs.
type ContextClientFunc func(context.Context) (*http.Client, error)

var contextClientFuncs []ContextClientFunc

func RegisterContextClientFunc(fn ContextClientFunc) {
	contextClientFuncs = append(contextClientFuncs, fn)
}

func ContextClient(ctx context.Context) (*http.Client, error) {
	if ctx != nil {
		if hc, ok := ctx.Value(HTTPClient).(*http.Client); ok {
			return hc, nil
		}
	}
	for _, fn := range contextClientFuncs {
		c, err := fn(ctx)
		if err != nil {
			return nil, err
		}
		if c != nil {
			return c, nil
		}
	}
	return http.DefaultClient, nil
}

func RetrieveToken(ctx context.Context, callback func(req *http.Request), tokenURL string, v url.Values) (*oauth2.Token, map[string]interface{}, error) {
	var token *oauth2.Token
	raw := map[string]interface{}{}
	hc, err := ContextClient(ctx)
	if err != nil {
		return token, raw, err
	}
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return token, raw, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if callback != nil {
		callback(req)
	}
	r, err := ctxhttp.Do(ctx, hc, req)
	if err != nil {
		return token, raw, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return token, raw, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}
	if code := r.StatusCode; code < 200 || code > 299 {
		return token, raw, fmt.Errorf("oauth2: cannot fetch token: %v\nResponse: %s", r.Status, body)
	}

	content, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch content {
	case "application/x-www-form-urlencoded", "text/plain":
		vals, err := url.ParseQuery(string(body))
		if err != nil {
			return token, raw, err
		}
		token = &oauth2.Token{
			AccessToken:  vals.Get("access_token"),
			TokenType:    vals.Get("token_type"),
			RefreshToken: vals.Get("refresh_token"),
		}

		for k, v := range vals {
			if len(v) > 0 {
				if len(v) == 1 {
					raw[k] = v[0]
				} else {
					raw[k] = v
				}
			}
		}

		e := vals.Get("expires_in")
		if len(e) == 0 {
			// TODO(jbd): Facebook's OAuth2 implementation is broken and
			// returns expires_in field in expires. Remove the fallback to expires,
			// when Facebook fixes their implementation.
			e = vals.Get("expires")
		}
		expires, _ := strconv.Atoi(e)
		if expires != 0 {
			token.Expiry = time.Now().Add(time.Duration(expires) * time.Second)
		}
	default:
		token = &oauth2.Token{}
		err = json.Unmarshal(body, &raw)
		if err != nil {
			return token, raw, err
		}
		var m map[string]interface{}
		if _, y := raw[`access_token`]; !y {
			for _, val := range raw {
				if v, y := val.(map[string]interface{}); y {
					if _, y := v[`access_token`]; y {
						m = v
						break
					}
				}
			}
		} else {
			m = raw
		}
		if m != nil {
			if v, y := m[`access_token`]; y {
				token.AccessToken = fmt.Sprint(v)
			}
			if v, y := m[`refresh_token`]; y {
				token.RefreshToken = fmt.Sprint(v)
			}
			if v, y := m[`token_type`]; y {
				token.TokenType = fmt.Sprint(v)
			}
			if v, y := m[`expires_in`]; y {
				lifetime, _ := strconv.ParseInt(fmt.Sprint(v), 10, 64)
				if lifetime > 0 {
					token.Expiry = time.Now().Add(time.Duration(lifetime) * time.Second)
				}
			}
			if v, y := m[`expires`]; y {
				lifetime, _ := strconv.ParseInt(fmt.Sprint(v), 10, 64)
				if lifetime > 0 {
					token.Expiry = time.Now().Add(time.Duration(lifetime) * time.Second)
				}
			}
		}
	}
	// Don't overwrite `RefreshToken` with an empty value
	// if this was a token refreshing request.
	if len(token.RefreshToken) == 0 {
		token.RefreshToken = v.Get("refresh_token")
	}
	return nil, raw, nil
}
