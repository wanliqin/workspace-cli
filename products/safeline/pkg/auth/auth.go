package auth

import "net/http"

// Transport adds authentication headers to requests.
type Transport struct {
	Base  http.RoundTripper
	Token string
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Token != "" {
		req.Header.Set("API-Token", t.Token)
	}
	return t.base().RoundTrip(req)
}

func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}
