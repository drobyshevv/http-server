package server

import "fmt"

// HandlerFunc handles an HTTP request.
type HandlerFunc func(r *Request, w ResponseWriter)

// Router routes HTTP requests to registered handlers.
type Router struct {
	Way map[string]HandlerFunc
}

// NewRouter creates a new router.
func NewRouter() *Router {
	return &Router{
		Way: make(map[string]HandlerFunc),
	}
}

// SetRoute registers a handler for an HTTP route.
func (rt *Router) SetRoute(method string, requestTarget string, handler HandlerFunc) {
	rt.Way[method+requestTarget] = handler
}

// Route routes an HTTP request to the corresponding handler.
func (rt *Router) Route(r *Request, w ResponseWriter) error {
	handler, found := rt.Way[r.Method+r.RequestTarget]
	if !found || handler == nil {
		return fmt.Errorf("route not found")
	}
	handler(r, w)
	return nil
}
