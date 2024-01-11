package lgo

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var DefaultServeMux = NewRouter()

// Handler http hander interface
type Handler interface {
	ServeHTTP(ctx *Context)
}

// WrapHandler MiddleWare For Wrap tgo.Handler To http.Handler
func WrapHandler(next Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := NewContext(w, r)

		defer func() {
			if r := recover(); r != nil {
				var err error
				switch r := r.(type) {
				case error:
					err = r
				default:
					err = fmt.Errorf("%v", r)
				}
				fmt.Println("err:", err)
			}
		}()

		defer ctx.WriteResponse()

		next.ServeHTTP(ctx)
	})
}

// HandlerFunc http hander function
type HandlerFunc func(ctx *Context)

func HandleFunc(pattern string, handler func(*Context)) *mux.Route {
	return DefaultServeMux.WrapHandle(pattern, HandlerFunc(handler))
}

func (f HandlerFunc) ServeHTTP(ctx *Context) {
	f(ctx)
}
