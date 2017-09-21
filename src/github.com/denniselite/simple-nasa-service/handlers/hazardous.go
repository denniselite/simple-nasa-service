package handlers

import (
	"github.com/valyala/fasthttp"
	"log"
	"encoding/json"
	s "github.com/denniselite/simple-nasa-service/service"
	"github.com/denniselite/simple-nasa-service/structs"
)

type HazardousHandlerResponse struct {
	Response []structs.NEO `json:"response"`
	Error string `json:"error"`
}

func (m *HazardousHandlerResponse) ToJSON() []byte {
	r, _ := json.Marshal(*m)
	return r
}

func HazardousHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("Application/JSON")
	res := new(HazardousHandlerResponse)
	var err error
	res.Response, err = s.GetNasaService().GetHazardous()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		res.Response = nil
		res.Error = err.Error()
		ctx.SetBody(res.ToJSON())
		log.Printf("Response: %+v\n", string(res.ToJSON()))
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(res.ToJSON())
	log.Printf("Response: %+v\n", string(res.ToJSON()))
	return
}
