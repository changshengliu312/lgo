package lgo

import (
	"github.com/gorilla/mux"
)

type Router struct {
	mux.Router
}

func NewRouter() *Router {
	//route := mux.NewRouter()
	return &Router{}
}

func (r *Router) WrapHandle(pattern string, handler Handler) *mux.Route {
	return r.Handle(pattern, WrapHandler(handler))
}
