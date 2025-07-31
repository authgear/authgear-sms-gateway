package sensitive

import (
	"errors"
	"net/url"
)

// RedactHTTPClientError redacts the error returned by http.Client.Do and friends.
// Note this function only works when err is *url.Error.
// This function does not work when *url.Error is wrapped.
//
// The documentation of http.Client.Do explicitly states that
// all errors returned by http.Client.Do is *url.Error.
// So this function does not need to handle wrapped *url.Error.
func RedactHTTPClientError(err error) error {
	var urlError *url.Error
	if errors.As(err, &urlError) {
		if urlError.URL != "" {
			if u, parseErr := url.Parse(urlError.URL); parseErr == nil {
				if u.RawQuery != "" {
					if q, parseErr := url.ParseQuery(u.RawQuery); parseErr == nil {
						for key := range q {
							q.Set(key, "REDACTED")
						}
						u.RawQuery = q.Encode()
						urlError.URL = u.String()
					}
				}
			}
		}
	}
	return err
}
