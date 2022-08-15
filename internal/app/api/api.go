package api

import (
	"context"
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

func (s *Server) Run(ctx context.Context) (err error) {
	port := os.Getenv("SERV_PORT")
	log.Println("Service will be started on", port)

	go func() {
		err = fasthttp.ListenAndServe(port, s.fastRouter.Handler)
	}()
	<-ctx.Done()
	return err
}

func (s *Server) configureRouter() {
	s.fastRouter.GET("/parameters", s.getParameters)
	s.fastRouter.GET("/health", s.health)
}
