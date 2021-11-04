// Package bunny provides functionality to interact with the Bunny CDN HTTP API.
package bunny

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/google/go-querystring/query"
)

const (
	// BaseURL is the base URL of the Bunny CDN HTTP API.
	BaseURL = "https://api.bunny.net"
	// AccessKeyHeaderKey is the name of the HTTP header that contains the Bunny API key.
	AccessKeyHeaderKey = "AccessKey"
	// DefaultUserAgent is the default value of the sent HTTP User-Agent header.
	DefaultUserAgent = "bunny-go"
)

// Logf is a log function signature.
type Logf func(format string, v ...interface{})

// Client is a Bunny CDN HTTP API Client.
type Client struct {
	baseURL *url.URL
	apiKey  string

	httpClient      http.Client
	httpRequestLogf Logf
	userAgent       string

	PullZone *PullZoneService
}

// NewClient returns a new bunny.net API client.
// The APIKey can be found in on the Account Settings page.
//
// Bunny.net API docs: https://support.bunny.net/hc/en-us/articles/360012168840-Where-do-I-find-my-API-key-
func NewClient(APIKey string, opts ...Option) *Client {
	clt := Client{
		baseURL:         mustParseURL(BaseURL),
		apiKey:          APIKey,
		httpClient:      *http.DefaultClient,
		userAgent:       DefaultUserAgent,
		httpRequestLogf: func(string, ...interface{}) {},
	}

	clt.PullZone = &PullZoneService{client: &clt}

	for _, opt := range opts {
		opt(&clt)
	}

	return &clt
}

func mustParseURL(urlStr string) *url.URL {
	res, err := url.Parse(urlStr)
	if err != nil {
		panic(fmt.Sprintf("Parsing url: %s failed: %s", urlStr, err))
	}

	return res
}

// newRequest creates an bunny.net API request.
// urlStr maybe absolute or relative, if it is relative it is joined with
// client.baseURL.
func (c *Client) newRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	url, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set(AccessKeyHeaderKey, c.apiKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// newGetRequest creates an bunny.NET API GET request.
// params must be a struct or nil, it is encoded into a query parameter.
// The struct must contain  `url` tags of the go-querystring package.
func (c *Client) newGetRequest(urlStr string, params interface{}) (*http.Request, error) {
	if params != nil {
		queryvals, err := query.Values(params)
		if err != nil {
			return nil, err
		}
		urlStr = urlStr + "?" + queryvals.Encode()
	}

	return c.newRequest(http.MethodGet, urlStr, nil)
}

// newPostRequest creates a bunny.NET API POST request.
// If body is not nil, it is encoded as JSON as send as HTTP-Body.
func (c *Client) newPostRequest(urlStr string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter

	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := c.newRequest(http.MethodPost, urlStr, buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) newDeleteRequest(urlStr string, params interface{}) (*http.Request, error) {
	return c.newRequest(http.MethodDelete, urlStr, nil)
}

// sendRequest sends a http Request to the bunny API.
// If the server returns a 2xx status code with an response body, the body is
// unmarshaled as JSON into result.
// If the ctx times out ctx.Error() is returned.
// If sending the response fails (http.Client.Do), the error will be returned.
// If the server returns an 401 error, an AuthenticationError error is returned.
// If the server returned an error and contains an APIError as JSON in the body,
// an APIError is returned.
// If the server returned a status code that is not 2xx an HTTPError is returned.
// If the HTTP request was successful, the response body is read and
// unmarshaled into result.
func (c *Client) sendRequest(ctx context.Context, req *http.Request, result interface{}) error {
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	c.logRequest(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok {
			if urlErr.Timeout() && ctx.Err() != nil {
				return ctx.Err()
			}
		}

		return err
	}

	defer resp.Body.Close() //nolint: errcheck

	if err := checkResp(req, resp); err != nil {
		return err
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return &HTTPError{
			RequestURL: req.URL.String(),
			StatusCode: resp.StatusCode,
			Errors:     []error{fmt.Errorf("reading response body failed: %w", err)},
		}
	}

	if result != nil {
		err = json.Unmarshal(buf, result)
		if err != nil {
			return &HTTPError{
				RequestURL: req.URL.String(),
				StatusCode: resp.StatusCode,
				RespBody:   buf,
				Errors:     []error{fmt.Errorf("decoding response body into json failed: %w", err)},
			}
		}
	}

	return nil
}

// checkResp checks if the resp indicates that the request was successful.
// If it wasn't an error is returned.
func checkResp(req *http.Request, resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			// ignore connection errors causing that the body can
			// not be received
			msg = []byte(http.StatusText(http.StatusUnauthorized))
		}

		return &AuthenticationError{
			Message: string(msg),
		}

	default:
		var err error

		httpErr := HTTPError{
			RequestURL: req.URL.String(),
			StatusCode: resp.StatusCode,
		}

		httpErr.RespBody, err = io.ReadAll(resp.Body)
		if err != nil {
			httpErr.Errors = append(httpErr.Errors, fmt.Errorf("reading response body failed: %w", err))

			return &httpErr
		}

		var apiErr APIError

		if err := json.Unmarshal(httpErr.RespBody, &apiErr); err != nil {
			httpErr.Errors = append(httpErr.Errors, fmt.Errorf("could not parse body as APIError: %w", err))
			return &httpErr
		}

		apiErr.HTTPError = httpErr
		return &apiErr
	}
}

func (c *Client) logRequest(req *http.Request) {
	if c.httpRequestLogf == nil {
		return
	}

	// hide the access key in the dumped request
	accessKey := req.Header.Get(AccessKeyHeaderKey)
	if accessKey != "" {
		req.Header.Set(AccessKeyHeaderKey, "***hidden***")
		defer func() { req.Header.Set(AccessKeyHeaderKey, accessKey) }()
	}

	debugReq, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		c.httpRequestLogf("dumping http request failed: %s", err)
		return
	}

	c.httpRequestLogf(string(debugReq))
}
