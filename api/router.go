package api

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Jidetireni/async-api.git/config"
	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	conf   *config.Config
	logger *slog.Logger
}

func New(conf *config.Config, logger *slog.Logger) *ApiServer {
	return &ApiServer{
		conf:   conf,
		logger: logger,
	}
}

func (s *ApiServer) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})

}

func (s *ApiServer) Start(ctx context.Context) error {
	router := gin.Default()
	router.GET("/ping", s.ping)

	// if err := router.Run(s.conf.ApiServerHost); err != nil {
	// 	s.logger.Error("Failed to run server", "error", err)
	// }

	srv := &http.Server{
		Addr:    net.JoinHostPort(s.conf.ApiServerHost, s.conf.ApiServerPort),
		Handler: router,
	}
	go func() {
		s.logger.Info("api server", "port", s.conf.ApiServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Failed to ListenAndServe:", "error", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownctx); err != nil {
			s.logger.Error("api server failed to shutdown", "error", err)
		}
	}()

	wg.Wait()
	return nil

}
