package handlers

import (
	"github.com/valyala/fasthttp"
	"log"
	"encoding/json"
)

type HomeHandlerResponse struct {
	Hello string `json:"hello"`
}

func (m *HomeHandlerResponse) ToJSON() []byte {
	r, _ := json.Marshal(*m)
	return r
}

func HomeHandler(ctx *fasthttp.RequestCtx) {
	res := new(HomeHandlerResponse)
	res.Hello = "world!"
	ctx.SetContentType("Application/JSON")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(res.ToJSON())
	log.Printf("Response: %+v\n", string(res.ToJSON()))
	return
}