package api

import (
	"context"
	"log/slog"
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

func NewApi(conf *config.Config, logger *slog.Logger) *ApiServer {
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
	middle := NewLoggerMiddleware(s.logger)
	router.Use(middle)

	// if err := router.Run(s.conf.ApiServerHost); err != nil {
	// 	s.logger.Error("Failed to run server", "error", err)
	// }

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		s.logger.Info("apiserver", "port", 8080)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Failed to ListenAndServe:", "error", err)
		}
	}()

	// Shutdown logic
	var wg sync.WaitGroup

	// This increments the WaitGroup counter by 1, indicating that one goroutine
	// is about to start.
	wg.Add(1)
	go func() {
		// ensures that wg.Done() is called when the goroutine exits
		// wg.Done() decrements the WaitGroup counter by 1
		defer wg.Done()
		<-ctx.Done() // used to signal cancellation or timeout. When the context is canceled

		//This creates a new context (shutdownctx) with a timeout of 10 seconds
		shutdownctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownctx); err != nil {
			s.logger.Error("api server failed to shutdown", "error", err)
		}
	}()

	wg.Wait()
	return nil

}
