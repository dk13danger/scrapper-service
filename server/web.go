package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/dk13danger/scrapper-service/config"
	"github.com/dk13danger/scrapper-service/service"
)

type Server struct {
	service *service.Service
	logger  *logrus.Logger
	cfg     *config.Server
}

func NewServer(
	service *service.Service,
	logger *logrus.Logger,
	cfg *config.Server,
) *Server {
	return &Server{
		service: service,
		logger:  logger,
		cfg:     cfg,
	}
}

func (s *Server) Run() {
	router := gin.Default()
	router.POST("/", s.handler())

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Port),
		Handler: router,
	}

	go func() {
		s.logger.Info("Starting server")
		srv.ListenAndServe()
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	s.logger.Println("Server shutdown started..")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.ShutdownTimeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Fatalf("Server shutdown err: %v", err)
	}
}
