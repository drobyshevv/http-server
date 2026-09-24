package server

import "fmt"

type HandlerFunc func(r *Request, w ResponseWriter)

type Router struct {
	Way map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		Way: make(map[string]HandlerFunc),
	}
}

func (rt *Router) SetRoute(method string, requestTarget string, handler HandlerFunc) {
	rt.Way[method+requestTarget] = handler
}

func (rt *Router) Route(r *Request, w ResponseWriter) error {
	handler, found := rt.Way[r.Method+r.RequestTarget]
	if !found || handler == nil {
		return fmt.Errorf("route not found")
	}
	handler(r, w)
	return nil
}
