package nhealth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/dghubble/oauth1"
	"github.com/google/go-querystring/query"
)

const (
	BaseURL = "https://developer.health.nokia.com"
)

type HealthAPI struct {
	oauthConfig oauth1.Config
}

var nhealth = oauth1.Endpoint{
	RequestTokenURL: "https://developer.health.nokia.com/account/request_token",
	AuthorizeURL:    "https://developer.health.nokia.com/account/authorize",
	AccessTokenURL:  "https://developer.health.nokia.com/account/access_token",
}

// NewHealthAPI creates a new Nokia Health API connector
func NewHealthAPI(apiKey, apiSecret, callback string) HealthAPI {
	config := oauth1.Config{
		ConsumerKey:            apiKey,
		ConsumerSecret:         apiSecret,
		CallbackURL:            callback,
		Endpoint:               nhealth,
		DisableCallbackConfirm: true,
	}

	return HealthAPI{oauthConfig: config}
}

// Request creates a new OAuth request instance
// Returns the request token, secret, and user auth URL on success
// Or an error on failure
func (hapi *HealthAPI) Request() (authToken, authSecret string, authURL *url.URL, err error) {
	authToken, authSecret, err = hapi.oauthConfig.RequestToken()
	if err != nil {
		return "", "", nil, err
	}

	authURL, err = hapi.oauthConfig.AuthorizationURL(authToken)
	if err != nil {
		return "", "", nil, err
	}

	return authToken, authSecret, authURL, nil
}

// Authorize handles an OAuth authorize response from the request url
// This returns a user access token and secret on success
func (hapi *HealthAPI) Authorize(token, secret string, req *http.Request) (userid, accessToken, accessSecret string, err error) {
	userid = req.FormValue("userid")

	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(req)
	if err != nil {
		return "", "", "", err
	}

	if token != requestToken {
		return "", "", "", fmt.Errorf("Mismatched request and response tokens")
	}

	accessToken, accessSecret, err = hapi.oauthConfig.AccessToken(token, secret, verifier)
	if err != nil {
		return "", "", "", err
	}

	return userid, accessToken, accessSecret, nil
}

// Get function wraps OAuth components to simplify other API calls
func (hapi *HealthAPI) Get(accessToken, accessSecret, url string, args interface{}) (*http.Response, error) {
	token := oauth1.NewToken(accessToken, accessSecret)
	config := oauth1.NewConfig(hapi.oauthConfig.ConsumerKey, hapi.oauthConfig.ConsumerSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Encode query
	v, err := query.Values(args)
	if err != nil {
		return nil, err
	}

	// Create request object
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = v.Encode()

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
