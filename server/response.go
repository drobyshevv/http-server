package server

import (
	"fmt"
	"net"
	"strings"
)

type Response struct {
	Protocol string
	Code     int
	Header   map[string]string
	Body     []byte
}

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

func (r *Response) SetStatus(code int) *Response {
	r.Code = code
	return r
}

func (r *Response) Status() int {
	return r.Code
}

func (r *Response) SetHeader(key string, val string) *Response {
	r.Header[key] = val
	return r
}

func (r *Response) SetBody(body []byte) *Response {
	r.Body = body
	r.Header["Content-Length"] = fmt.Sprintf("%d", len(body))
	return r
}

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
