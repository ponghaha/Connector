package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new Zap SugaredLogger with a pre-defined enterprise-friendly configuration.
func New(level string) *zap.SugaredLogger {
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		// Default to InfoLevel if the provided level is invalid.
		logLevel = zapcore.InfoLevel
	}

	// A manually crafted config gives us more control over the final log structure.
	config := zap.Config{
		Encoding: "json",
		Level:    zap.NewAtomicLevelAt(logLevel),
		// Corrected syntax: removed extra curly braces
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "timestamp",
			CallerKey:    "caller",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder, // Standard time format
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, _ := config.Build()
	return logger.Sugar()
}

// GinLogger is a Gin middleware that logs requests using our configured SugaredLogger.
func GinLogger(logger *zap.SugaredLogger, apiRequestID string, apiLanguage string, apiDeviceOS string, apiChannel string,) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process the request
		c.Next()

		latency := time.Since(start)

		// Check for errors written to the context
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.Errorw("Request error",
					"error", e.Err,
					"apiRequestID", c.GetString(apiRequestID),
					"apiLanguage", c.GetString(apiLanguage),
					"apiDeviceOS", c.GetString(apiDeviceOS),
					"apiChannel", c.GetString(apiChannel),
				)
			}
		} else {
			// Log successful requests
			logger.Infow("Request handled",
				"status", c.Writer.Status(),
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"query", c.Request.URL.RawQuery,
				"ip", c.ClientIP(),
				"latency", latency.String(),
				"user_agent", c.Request.UserAgent(),
				"apiRequestID", c.GetString(apiRequestID),
				"apiLanguage", c.GetString(apiLanguage),
				"apiDeviceOS", c.GetString(apiDeviceOS),
				"apiChannel", c.GetString(apiChannel),
			)
		}
	}
}
