package handlers

import (
	"github.com/valyala/fasthttp"
	"log"
	"encoding/json"
)

type MainHandlerRequest struct{}

type MainHandlerResponse struct {
	Hello string `json:"hello"`
}

func (m *MainHandlerResponse) ToJSON() []byte {
	r, _ := json.Marshal(*m)
	return r
}

func MainHandler(ctx *fasthttp.RequestCtx) {
	res := new(MainHandlerResponse)
	res.Hello = "world!"
	log.Printf("Response: %+v\n", *res)
	ctx.SetContentType("Application/JSON")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(res.ToJSON())
}