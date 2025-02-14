package main

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func setupServer() *gin.Engine {
	ginApp := gin.Default()

	log := SetupLogger()
	ginApp.Use(LoggerMiddleware(log))

	log.Info("starting server", slog.String("host", os.Getenv("APP_HOST")+":"+os.Getenv("APP_PORT")))
	SetupRoutes(ginApp)

	return ginApp
}

func SetupAndRunServer() {
	err := setupServer().Run(":" + os.Getenv("APP_PORT"))

	if err != nil {
		panic("can't start server " + err.Error())
	}
}

func LoggerMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		log.Info("request started",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("ip", c.ClientIP()),
			slog.String("user-agent", c.Request.UserAgent()),
		)

		rw := &responseWriter{c.Writer, http.StatusOK}
		c.Writer = rw

		c.Next()

		log.Info("request completed",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", rw.status),
			slog.Duration("duration", time.Since(start)),
		)
	}
}

type responseWriter struct {
	gin.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
