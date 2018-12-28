package httpServer

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/valyala/fasthttp"
	"time"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

func GenerateResp(result string, code int, message string) map[string]map[string]interface{} {
	resp := make(map[string]map[string]interface{})
	resp["data"] = make(map[string]interface{})
	resp["data"]["result"] = result
	resp["status"] = make(map[string]interface{})
	resp["status"]["code"] = code
	resp["status"]["message"] = message
	return resp
}

func DoJSONWrite(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	start := time.Now()
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		elapsed := time.Since(start)
		logs.Error("", elapsed, err.Error(), obj)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}
