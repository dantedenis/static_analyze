package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"net"
	"static-analyze/pkg/indicator"
	"static-analyze/pkg/indicator/mock"
	"testing"
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

	if resp.StatusCode() != 200 {
		t.Error("Error status code")
	}
	if string(resp.Body()) != "{\"message\":\"200\"}" {
		t.Error("Error Body response", string(resp.Body()))
	}
}

func TestNewServer(t *testing.T) {
	assert.NotNil(t, NewServer(nil))
}

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
