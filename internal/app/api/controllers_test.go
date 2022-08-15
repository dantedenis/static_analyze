package api

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"net"
	"net/http"
	"static-analyze/pkg/indicator"
	"static-analyze/pkg/indicator/mock"
	"testing"
	"time"
)

var (
	pairStr = []string{"VALUE1", "VALUE3", "", "TEST"}
)

func Test_Health(t *testing.T) {
	port := ":1234"
	host := "http://localhost"
	serv := NewServer(indicator.New(mock.Proto{}, pairStr))

	defer startServerOnPort(t, port, serv.health).Close()

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(host + port + "/health")
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")

	resp := fasthttp.AcquireResponse()

	assert.Nil(t, fasthttp.Do(req, resp))

	if resp.StatusCode() != http.StatusOK {
		t.Error("Error status code")
	}
	if string(resp.Body()) != `{"message":"200"}` {
		t.Error("Error Body response", string(resp.Body()))
	}
}

////////////////////////////////////////////////////////////////////////////////////////

var (
	reqBody = `
	{
	    "currency":"VALUE1",
	    "period": 3600
	}`
	respBody = `{"last_update":"0001-01-01T00:00:00Z","currency_pair":"VALUE1","timeframe":"60 min","indicators":[]}`
)

func Test_GetParameters(t *testing.T) {
	port := ":1235"
	host := "http://localhost"
	serv := NewServer(indicator.New(mock.Proto{}, pairStr))

	defer startServerOnPort(t, port, serv.getParameters).Close()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(host + port + "/parameters")
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")

	req.SetBody([]byte(reqBody))

	resp := fasthttp.AcquireResponse()

	assert.Nil(t, fasthttp.Do(req, resp))

	if resp.StatusCode() != http.StatusOK {
		t.Error("Error status code")
	}
	if string(resp.Body()) != respBody {
		t.Error("error body response")
	}
}

func Test_GetParameters_error(t *testing.T) {
	port := ":1236"
	host := "http://localhost"
	serv := NewServer(indicator.New(mock.Proto{}, pairStr))

	defer startServerOnPort(t, port, serv.getParameters).Close()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(host + port + "/parameters")
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")

	resp := fasthttp.AcquireResponse()

	assert.Nil(t, fasthttp.Do(req, resp))

	if resp.StatusCode() != http.StatusBadRequest {
		t.Error("Error status code")
	}
}

func Test_GetParameters_error2(t *testing.T) {
	port := ":1237"
	host := "http://localhost"
	serv := NewServer(indicator.New(mock.Proto{}, pairStr))

	defer startServerOnPort(t, port, serv.getParameters).Close()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(host + port + "/parameters")
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")

	resp := fasthttp.AcquireResponse()

	assert.Nil(t, fasthttp.Do(req, resp))

	if resp.StatusCode() != http.StatusBadRequest {
		t.Error("Error status code")
	}
}

////////////////////////////////////////////////////////////////////////////////////////

var (
	reqBodyFalsePeriod = `
	{
	    "currency":"VALUE1",
	    "period": 360
	}`
	reqBodyFalseCurrency = `
	{
	    "currency":"VA",
	    "period": 3600
	}`
	respBodyFalsePeriod   = `{"message":"error period"}`
	respBodyFalseCurrency = `{"message":"does not exist key-currency: VA"}`
)

func Test_Constraint_Period(t *testing.T) {
	port := ":1238"
	host := "http://localhost"
	serv := NewServer(indicator.New(mock.Proto{}, pairStr))

	defer startServerOnPort(t, port, serv.getParameters).Close()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(host + port + "/parameters")
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")

	resp := fasthttp.AcquireResponse()

	req.SetBody([]byte(reqBodyFalsePeriod))

	assert.Nil(t, fasthttp.Do(req, resp))

	if string(resp.Body()) != respBodyFalsePeriod {
		t.Error("Uncorrected body")
	}

	if resp.StatusCode() != http.StatusBadRequest {
		t.Error("Uncorrected status")
	}
}

func Test_Constraint_Currency(t *testing.T) {
	port := ":1238"
	host := "http://localhost"
	serv := NewServer(indicator.New(mock.Proto{}, pairStr))

	defer startServerOnPort(t, port, serv.getParameters).Close()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(host + port + "/parameters")
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")

	resp := fasthttp.AcquireResponse()

	req.SetBody([]byte(reqBodyFalseCurrency))

	assert.Nil(t, fasthttp.Do(req, resp))

	if string(resp.Body()) != respBodyFalseCurrency {
		t.Error("Uncorrected body")
	}

	if resp.StatusCode() != http.StatusBadRequest {
		t.Error("Uncorrected status")
	}
}

////////////////////////////////////////////////////////////////////////////////////////

func startServerOnPort(t *testing.T, port string, h fasthttp.RequestHandler) io.Closer {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost%s", port))
	if err != nil {
		t.Fatalf("cannot start tcp server on port %s: %s", port, err)
	}

	go func() {
		err = fasthttp.Serve(ln, h)
		if err != nil {
			log.Println(err)
		}
	}()

	return ln
}

func TestServer_Run(t *testing.T) {
	serv := NewServer(nil)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(500*time.Millisecond))
	defer cancel()
	assert.Nil(t, serv.Run(ctx))
}
