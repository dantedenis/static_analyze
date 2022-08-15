package api

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	indic "static-analyze/pkg/indicator"
)

type Server struct {
	fastRouter *router.Router
	service    indic.Get
}

func NewServer(serv indic.Get) *Server {
	s := &Server{
		fastRouter: router.New(),
		service:    serv,
	}

	s.configureRouter()

	return s
}

func (s *Server) Run() error {
	port := os.Getenv("SERV_PORT")
	log.Println("Service will be started on", port)

	return fasthttp.ListenAndServe(port, s.fastRouter.Handler)
}

func (s *Server) configureRouter() {
	s.fastRouter.GET("/parameters", s.getParameters)
	s.fastRouter.GET("/health", s.health)
}
