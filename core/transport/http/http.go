package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*http.Server
	opts options
}

func NewHTTPServer(opts ...Options) *Server {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	hs := &Server{
		opts: o,
	}

	r := gin.Default()
	r.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World, this is Gin admin"})
	})
	r.GET("/account", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World, this is Gin account."})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", o.Port),
		Handler: r,
	}
	hs.Server = srv

	return hs
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen %s\n", err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.Shutdown(ctx); err != nil {
		fmt.Printf("http server shutdown failed:%s\n", err)
		return err
	}
	fmt.Printf("http server shutdown success\n")
	return nil
}
