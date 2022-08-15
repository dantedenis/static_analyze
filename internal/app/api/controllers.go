package api

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	indic "static-analyze/pkg/indicator"
	"time"
)

var (
	timeframe = []int64{60, 300, 900, 1800, 3600}
)

// request object
type request struct {
	Cur    string `json:"currency"`
	Period int64  `json:"period"`
}

// response object
type response struct {
	LastUpdate time.Time         `json:"last_update"`
	Currency   string            `json:"currency_pair"`
	Period     string            `json:"timeframe"`
	Indicator  []indic.Indicator `json:"indicators"`
}

func (s *Server) health(ctx *fasthttp.RequestCtx) {

	JSONResponse(ctx, http.StatusOK, map[string]string{
		"message": "200",
	})
}

func (s *Server) getParameters(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	req := request{}
	err := json.Unmarshal(body, &req)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	if !containTimeframe(req.Period) {
		JSONResponse(ctx, http.StatusBadRequest, map[string]string{
			"message": "error period",
		})
		return
	}

	val, timeUpdate, err := s.service.Get(req.Cur, req.Period)
	if err != nil {
		JSONResponse(ctx, http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	JSONResponse(ctx, http.StatusOK, response{
		Currency:   req.Cur,
		Period:     fmt.Sprintf("%d min", req.Period/60),
		Indicator:  val,
		LastUpdate: timeUpdate,
	})
}

func JSONResponse(ctx *fasthttp.RequestCtx, code int, data interface{}) {
	b, _ := json.Marshal(data)

	ctx.Response.Header.Set("Content-Type", "application/json; charset=utf-8")
	ctx.SetStatusCode(code)

	_, err := ctx.Write(b)
	if err != nil {
		log.Println(err)
		return
	}
}

func containTimeframe(target int64) bool {
	for _, k := range timeframe {
		if target == k {
			return true
		}
	}
	return false
}
