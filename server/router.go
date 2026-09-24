package server

import "fmt"

// HandlerFunc обрабатывает HTTP-запрос.
type HandlerFunc func(r *Request, w ResponseWriter)

// Router маршрутизирует HTTP-запросы к зарегистрированным обработчикам.
type Router struct {
	Way map[string]HandlerFunc
}

// NewRouter создаёт новый роутер.
func NewRouter() *Router {
	return &Router{
		Way: make(map[string]HandlerFunc),
	}
}

// SetRoute регистрирует обработчик для HTTP-маршрута.
func (rt *Router) SetRoute(method string, requestTarget string, handler HandlerFunc) {
	rt.Way[method+requestTarget] = handler
}

// Route маршрутизирует HTTP-запрос к соответствующему обработчику.
func (rt *Router) Route(r *Request, w ResponseWriter) error {
	handler, found := rt.Way[r.Method+r.RequestTarget]
	if !found || handler == nil {
		return fmt.Errorf("route not found")
	}
	handler(r, w)
	return nil
}
