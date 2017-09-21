package handlers

import (
	"github.com/valyala/fasthttp"
	"log"
	"encoding/json"
	s "github.com/denniselite/simple-nasa-service/service"
	"github.com/denniselite/simple-nasa-service/structs"
	"strings"
)

const (
	isHazardousTrue = "true"
	isHazardousFalse = "false"
)

type FastestHandlerResponse struct {
	Response *structs.NEO `json:"response"`
	Error string `json:"error"`
}

func (m *FastestHandlerResponse) ToJSON() []byte {
	r, _ := json.Marshal(*m)
	return r
}

func FastestHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("Application/JSON")
	res := new(FastestHandlerResponse)

	if !ctx.QueryArgs().Has("is_hazardous") {
		res.Error = `"is_hazardous"" option is required`
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody(res.ToJSON())
		log.Printf("Response: %+v\n", string(res.ToJSON()))
		return
	}

	isHazardousStr := strings.ToLower(string(ctx.QueryArgs().Peek("is_hazardous")))

	if (isHazardousStr != isHazardousTrue) && (isHazardousStr != isHazardousFalse) {
		res.Error = `"is_hazardous"" option should contains only "true" or "false"`
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody(res.ToJSON())
		log.Printf("Response: %+v\n", string(res.ToJSON()))
		return
	}

	isHazardous := isHazardousStr == isHazardousTrue
	fastest, err := s.GetNasaService().GetFastest(isHazardous)
	res.Response = &fastest
	if err != nil {
		res.Response = nil
		res.Error = err.Error()
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody(res.ToJSON())
		log.Printf("Response: %+v\n", string(res.ToJSON()))
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(res.ToJSON())
	log.Printf("Response: %+v\n", string(res.ToJSON()))
	return
}
