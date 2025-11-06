package handler

import (
	"encoding/json"
	"log"

	"github.com/benlocal/lai-panel/pkg/di"
	"github.com/valyala/fasthttp"
)

func HandleWithDI(fn interface{}) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		err := di.Invoke(fn)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetBodyString(err.Error())
		}
	}
}

func RouteDI(handlerFn interface{}) fasthttp.RequestHandler {
	return HandleWithDI(handlerFn)
}

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSONSuccess(ctx *fasthttp.RequestCtx, data interface{}) {
	resp := response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(resp)
}

func JSONError(ctx *fasthttp.RequestCtx, message string, err error) {
	if err != nil {
		log.Printf("Failed with error: %v", err)
	}

	resp := response{
		Code:    -1,
		Message: message,
		Data:    nil,
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")

	json.NewEncoder(ctx).Encode(resp)
}
