package twitch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/atasik/go-twitch-sdk/internal/request"
	"github.com/atasik/go-twitch-sdk/pkg/input"
	"github.com/atasik/go-twitch-sdk/pkg/response"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/google/go-querystring/query"

	"github.com/pkg/errors"
)

const (
	host           = "https://api.twitch.tv/helix"
	authorizeUrl   = "https://id.twitch.tv/oauth2/authorize?response_type=%s&client_id=%s&redirect_uri=%s&scope=%s&state=%s"
	authorizeToken = "https://id.twitch.tv/oauth2/token"

	endpointUsers         = "/users"
	endpointSubscriptions = "/eventsub/subscriptions"

	defaultTimeout = 5 * time.Second
)

const (
	ChannelFollow = "channel.follow"
	StreamOnline  = "stream.online"
	StreamOffline = "stream.offline"
)

// Client is a getpocket API client
type Client struct {
	client       *http.Client
	clientId     string
	clientSecret string
}

// NewClient creates a new client instance with your client id and your secret code
func NewClient(clientId, secredCode string) (*Client, error) {
	if clientId == "" || secredCode == "" {
		return nil, errors.New("empty params")
	}

	return &Client{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		clientId:     clientId,
		clientSecret: secredCode,
	}, nil
}

// GetAuthorizationURL generates link to authorize user
func (c *Client) GetAuthorizationURL(redirectUri, state, scope, responseType string) (string, error) {
	if redirectUri == "" {
		return "", errors.New("empty params")
	}

	return fmt.Sprintf(authorizeUrl, responseType, c.clientId, redirectUri, scope, state), nil
}

// GetAccessToken obtains the access token that is used to do request to TwitchAPI
func (c *Client) GetAccessToken(ctx context.Context, code, gt, redirectUri string) (*response.AuthorizeResponse, error) {
	inp := &request.AuthorizeRequest{
		ClientId:     c.clientId,
		ClientSecret: c.clientSecret,
		GrantType:    gt,
		RedirectUri:  redirectUri,
		Code:         code,
	}

	v, _ := query.Values(inp)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authorizeToken, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create new request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.responseHandler(req)
	if err != nil {
		return nil, err
	}

	var authResp response.AuthorizeResponse
	json.Unmarshal([]byte(resp), &authResp)

	caser := cases.Title(language.English)
	authResp.TokenType = caser.String(authResp.TokenType)

	if authResp.AccessToken == "" {
		return nil, errors.New("empty access token in API response")
	}

	return &authResp, nil
}

// GetUser gives infromation about user
func (c *Client) GetUser(ctx context.Context, inp input.UserInput, accessToken string) (*response.UserResponse, error) {
	resp, err := c.doHTTP(ctx, inp, nil, endpointUsers, accessToken, http.MethodGet)

	var userResp response.UserResponse
	json.Unmarshal([]byte(resp), &userResp)

	return &userResp, err
}

// Subscribe is used to subscribe to some events
func (c *Client) Subscribe(ctx context.Context, typ, id, callback, secret, accessToken string) (*response.SubResponse, error) {
	inp := &request.SubRequest{
		Type:    typ,
		Version: "1",
		Condition: request.Condition{
			Id: id,
		},
		Transport: request.Transport{
			Method:   "webhook",
			Callback: callback,
			Secret:   secret,
		},
	}
	resp, err := c.doHTTP(ctx, nil, inp, endpointSubscriptions, accessToken, http.MethodPost)

	var subResp response.SubResponse
	json.Unmarshal([]byte(resp), &subResp)

	return &subResp, err
}

// DeleteSubscriptions is used to delete subscription
func (c *Client) DeleteSubscription(ctx context.Context, inp input.UserInput, accessToken string) error {
	if _, err := c.doHTTP(ctx, inp, nil, endpointSubscriptions, accessToken, http.MethodDelete); err != nil {
		return err
	}
	return nil
}

// GetSubscriptions is used to know about event you currently subscribed to
func (c *Client) GetSubscriptions(ctx context.Context, accessToken string) (*response.SubResponse, error) {
	respB, err := c.doHTTP(ctx, nil, nil, endpointSubscriptions, accessToken, http.MethodGet)

	var subResp response.SubResponse
	json.Unmarshal([]byte(respB), &subResp)
	return &subResp, err
}

func (c *Client) doHTTP(ctx context.Context, params, body interface{}, endpoint, accessToken, method string) (string, error) {
	v, _ := query.Values(params)
	b, _ := json.Marshal(body)

	fmt.Println(string(b))

	req, err := http.NewRequestWithContext(ctx, method, host+endpoint+fmt.Sprintf("?%s", v.Encode()), bytes.NewBuffer(b))
	if err != nil {
		return "", errors.WithMessage(err, "failed to create new request")
	}

	req.Header.Set("Authorization", accessToken)
	req.Header.Set("Client-Id", c.clientId)
	req.Header.Set("Content-Type", "application/json")

	return c.responseHandler(req)
}

//handler
func (c *Client) responseHandler(req *http.Request) (string, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return "", errors.WithMessage(err, "failed to send http request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusNoContent {
		err := fmt.Sprintf("API Error: %d", resp.StatusCode)
		return "", errors.New(err)
	}

	respB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.WithMessage(err, "failed to read request params")
	}

	return string(respB), nil
}
