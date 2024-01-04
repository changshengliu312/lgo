package lgo

import (
	"context"
	"lgo/utils"
	"net/http"
)

//思考:?需要传递哪些上下文内容
type Context struct {
	context.Context
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	ClientIP       uint32
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{
		Context: context.WithValue(r.Context(), "retRode", 0),

		Request:        r,
		ResponseWriter: w,
	}
	ctx.ClientIP = utils.AddrtoI(ctx.Request.RemoteAddr)
	return ctx
}
