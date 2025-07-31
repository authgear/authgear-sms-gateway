package sensitive

import (
	"errors"
	"fmt"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRedactHTTPClientError(t *testing.T) {
	Convey("RedactHTTPClientError", t, func() {
		Convey("should handle nil input", func() {
			result := RedactHTTPClientError(nil)
			So(result, ShouldBeNil)
		})

		Convey("should return non-url.Error unchanged", func() {
			err := errors.New("generic error")
			result := RedactHTTPClientError(err)
			So(result.Error(), ShouldEqual, "generic error")
		})

		Convey("should handle url.Error with empty URL", func() {
			urlErr := &url.Error{
				Op:  "Get",
				URL: "",
				Err: errors.New("connection failed"),
			}
			result := RedactHTTPClientError(urlErr)
			So(result.Error(), ShouldEqual, "Get \"\": connection failed")
		})

		Convey("should handle url.Error with invalid URL", func() {
			urlErr := &url.Error{
				Op:  "Get",
				URL: "://invalid-url",
				Err: errors.New("connection failed"),
			}
			result := RedactHTTPClientError(urlErr)
			So(result.Error(), ShouldEqual, "Get \"://invalid-url\": connection failed")
		})

		Convey("should handle url.Error with URL without query", func() {
			urlErr := &url.Error{
				Op:  "Get",
				URL: "https://example.com/path",
				Err: errors.New("connection failed"),
			}
			result := RedactHTTPClientError(urlErr)
			So(result.Error(), ShouldEqual, "Get \"https://example.com/path\": connection failed")
		})

		Convey("should redact query parameters in url.Error", func() {
			urlErr := &url.Error{
				Op:  "Get",
				URL: "https://example.com/path?api_key=secret123&token=abc456&user=john",
				Err: errors.New("connection failed"),
			}
			result := RedactHTTPClientError(urlErr)
			So(result.Error(), ShouldEqual, "Get \"https://example.com/path?api_key=REDACTED&token=REDACTED&user=REDACTED\": connection failed")
		})

		Convey("should handle url.Error with malformed query", func() {
			urlErr := &url.Error{
				Op:  "Get",
				URL: "https://example.com/path?%zzz",
				Err: errors.New("connection failed"),
			}
			result := RedactHTTPClientError(urlErr)
			So(result.Error(), ShouldEqual, "Get \"https://example.com/path?%zzz\": connection failed")
		})

		Convey("should preserve URL structure while redacting", func() {
			urlErr := &url.Error{
				Op:  "Post",
				URL: "https://api.example.com:8080/v1/sms?key=secret&msg=hello%20world",
				Err: errors.New("timeout"),
			}
			result := RedactHTTPClientError(urlErr)
			So(result.Error(), ShouldEqual, "Post \"https://api.example.com:8080/v1/sms?key=REDACTED&msg=REDACTED\": timeout")
		})

		Convey("should handle empty query values", func() {
			urlErr := &url.Error{
				Op:  "Get",
				URL: "https://example.com/path?empty=&key=value",
				Err: errors.New("error"),
			}
			result := RedactHTTPClientError(urlErr)
			So(result.Error(), ShouldEqual, "Get \"https://example.com/path?empty=REDACTED&key=REDACTED\": error")
		})

		Convey("does not handle wrapped *url.URL", func() {
			urlErr := &url.Error{
				Op:  "Post",
				URL: "https://api.example.com/send?token=secret123&message=hello",
				Err: errors.New("network error"),
			}
			wrappedErr := fmt.Errorf("failed to send: %w", urlErr)
			result := RedactHTTPClientError(wrappedErr)
			So(result.Error(), ShouldEqual, "failed to send: Post \"https://api.example.com/send?token=secret123&message=hello\": network error")
		})
	})
}
