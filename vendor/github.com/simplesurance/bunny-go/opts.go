package bunny

// Option is a type for Client options.
type Option func(*Client)

// WithHTTPRequestLogger is an option to log all sent out HTTP-Request via a log function.
func WithHTTPRequestLogger(logger Logf) Option {
	return func(clt *Client) {
		clt.httpRequestLogf = logger
	}
}

// WithUserAgent is an option to specify the value of the User-Agent HTTP
// Header.
func WithUserAgent(userAgent string) Option {
	return func(clt *Client) {
		clt.userAgent = userAgent
	}
}
