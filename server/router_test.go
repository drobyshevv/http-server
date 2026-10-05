package server

import (
	"testing"
)

func TestCalledRouter(t *testing.T) {
	called := false
	handler := func(r *Request, w ResponseWriter) {
		called = true
	}

	r := NewRouter()
	r.SetRoute("GET", "/users", handler)

	req := &Request{
		Method:        "GET",
		RequestTarget: "/users",
	}

	err := r.Route(req, nil)

	if !called {
		t.Error("handler was not called")
	}

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestRouteNotFound(t *testing.T) {
	rt := NewRouter()

	handler := func(r *Request, w ResponseWriter) {}

	rt.SetRoute("GET", "/users", handler)

	req := &Request{
		Method:        "GET",
		RequestTarget: "/posts",
	}

	err := rt.Route(req, nil)

	if err == nil {
		t.Error("expected error, got nil")
	}
}
