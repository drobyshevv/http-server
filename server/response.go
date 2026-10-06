package server

import (
	"fmt"
	"net"
	"strings"
)

// Response represents an HTTP response.
type Response struct {
	Protocol string
	Code     int
	Header   map[string]string
	Body     []byte
}

// NewResponse creates a new HTTP response with the specified protocol
// and a 200 OK status.
func NewResponse(
	protocol string,
) *Response {
	return &Response{
		Protocol: protocol,
		Code:     200,
		Header: map[string]string{
			"Content-Length": "0",
		},
	}
}

// ResponseWriter defines methods for constructing an HTTP response.
type ResponseWriter interface {
	SetStatus(int)
	Status() int
	SetHeader(string, string)
	SetBody([]byte)
}

// SetStatus sets the HTTP response status code.
func (r *Response) SetStatus(code int) {
	r.Code = code
}

// Status returns the HTTP response status code.
func (r *Response) Status() int {
	return r.Code
}

// SetHeader sets an HTTP response header.
func (r *Response) SetHeader(key string, val string) {
	r.Header[key] = val
}

// SetBody sets the response body and updates the Content-Length header.
func (r *Response) SetBody(body []byte) {
	r.Body = body
	r.Header["Content-Length"] = fmt.Sprintf("%d", len(body))
}

// write writes the HTTP response to a TCP connection.
func (r *Response) write(conn net.Conn) error {
	statusLine := fmt.Sprintf("%s %d %s\r\n", r.Protocol, r.Code, statusText(r.Code))

	var headerLines strings.Builder
	for key, val := range r.Header {
		headerLines.WriteString(fmt.Sprintf("%s: %s\r\n", key, val))
	}
	headerLines.WriteString("\r\n")

	_, err := conn.Write([]byte(statusLine + headerLines.String()))
	if err != nil {
		return err
	}

	_, err = conn.Write(r.Body)
	if err != nil {
		return err
	}

	return nil
}

func statusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 204:
		return "No Content"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return ""
	}
}
