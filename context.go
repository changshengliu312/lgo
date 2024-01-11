package lgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"lgo/retCode"
	"lgo/utils"
	"net/http"
)

// retcode
var (
	RetCodeKey = "retcode"
	RetMsgKey  = "retmsg"
	ResultKey  = "result"
)

//response type 回报报文格式
const (
	JSONType     = 1
	RawBytesType = 2
)

//思考:?需要传递哪些上下文内容
type Context struct {
	context.Context
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	ContentType    string
	ResponseType   int

	Result     map[string]interface{} // for json
	ResultBody interface{}            //回报数据
	RawRspBody []byte                 // for raw bytes
	ClientIP   uint32

	HTMLEscape bool //转换json时 是否转义
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{
		Context: context.WithValue(r.Context(), "retRode", 0),

		Request:        r,
		ResponseWriter: w,
		ResponseType:   JSONType,
		Result:         make(map[string]interface{}),
		HTMLEscape:     true,
	}
	ctx.ClientIP = utils.AddrtoI(ctx.Request.RemoteAddr)
	return ctx
}

// SetHeader set header
func (ctx *Context) SetHeader(key, val string) {
	ctx.ResponseWriter.Header().Set(key, val)
}

// GetHeader get header
func (ctx *Context) GetHeader(key string) string {
	return ctx.Request.Header.Get(key)
}

func (ctx *Context) WriteResponse() {
	if ctx.ResponseType == JSONType {
		ctx.WriteJSON()
	} else {
		ctx.WriteRawBytes()
	}
}

func (ctx *Context) SetResultBody(res interface{}) {
	ctx.ResultBody = res
	ctx.Result[ResultKey] = res
}

// SetBytesResult set response body for raw bytes
func (ctx *Context) SetBytesResult(data []byte) error {
	//ctx.ResponseType = RspTypeRawBytes
	ctx.RawRspBody = data
	return nil
}

// WriteRawBytes 返回原始数据流
func (ctx *Context) WriteRawBytes() {
	if ctx.ContentType == "" {
		ctx.SetHeader("Content-Type", "application/octet-stream; charset=utf-8")
	}
	ctx.ResponseWriter.Write(ctx.RawRspBody)
}

// JSONMarshal json marshal with encoder
func (ctx *Context) JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(ctx.HTMLEscape) /*字符串在编码为JSON字符串时会被强制转换为有效的UTF-8，为了防止一些浏览器在JSON输出误解以为是HTML，“<”，“>”，“&”这类字符会被进行转义*/
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// WriteJSON 返回json格式 默认
func (ctx *Context) WriteJSON() {
	if ctx.ContentType == "" {
		ctx.SetHeader("Content-Type", "application/json; charset=utf-8")
	}
	if ret := ctx.Value("lgoRetcode"); ret != nil {
		if res, ok := ret.(*retCode.Retcode); ok {
			ctx.Result[RetCodeKey] = res.ErrCode
			ctx.Result[RetMsgKey] = res.ErrMsg
		}
	}

	if len(ctx.Result) > 0 {
		data, err := ctx.JSONMarshal(ctx.Result)
		if err != nil {
			fmt.Printf("lgo json marshal fail:%v\n", err)
		}
		//dataLen := len(data)
		ctx.ResponseWriter.Write(data)
	}
}
