package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"connectorapi-go/internal/adapter/utils"
	"connectorapi-go/pkg/logger"
	"connectorapi-go/pkg/metrics"
	_ "connectorapi-go/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

const apiRequestID = "Api-RequestID"
const apiKey       = "Api-Key"
const apiLanguage  = "Api-Language"
const apiDeviceOS  = "Api-DeviceOS"
const apiChannel   = "Api-Channel"

// SetupRouter
func SetupRouter(
	appLogger *zap.SugaredLogger,
	repo *utils.APIKeyRepository,
	collectionHandler *collectionHandler,
	agreementHandler *agreementHandler,
	creditCardHandler *creditCardHandler,
) *gin.Engine {
	router := gin.New()

	// --- Global Middlewares ---
	router.Use(ApiRequestIDMiddleware())
	router.Use(ApiKeyMiddleware())
	router.Use(ApiLanguageMiddleware())
	router.Use(ApiDeviceOSMiddleware())
	router.Use(ApiChannelMiddleware())

	router.Use(logger.GinLogger(appLogger, apiRequestID, apiLanguage, apiDeviceOS, apiChannel))
	router.Use(PrometheusMiddleware())
	router.Use(gin.Recovery())

	// --- Public API Group ---

	router.GET("/healthz", HealthCheck)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- API Group  ---
	apiRoute := router.Group("/Api")
	{
		collectionHandler.RegisterRoutes(apiRoute)
		agreementHandler.RegisterRoutes(apiRoute)
		creditCardHandler.RegisterRoutes(apiRoute)
	}

	return router
}

// --- Middlewares Definitions ---
// RequestIDMiddleware checks for an incoming X-Request-ID, RequestID header
func ApiRequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := utils.GetHeader(c, "X-Request-ID", "Api-RequestID", "RequestID")

		if reqID == "" {
			prefix := "RQ"
			now := time.Now()
			dateTime := now.Format("20060102150405")                // YYYYMMDDhhmmss

			r := rand.New(rand.NewSource(time.Now().UnixNano()))    // seed random with time now
			runningNo := r.Intn(10000)                              // random 0–9999
			runningStr := fmt.Sprintf("%04d", runningNo)

			reqID = prefix + dateTime + runningStr
		}

		c.Set(apiRequestID, reqID)
		c.Header("X-Request-ID", reqID)
		c.Header("Api-RequestID", reqID)

		c.Next()
	}
}

func ApiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := utils.GetHeader(c, "X-Key", "Api-Key", "APIKey")

		c.Set(apiKey, key)
		c.Header("X-Key", key)
		c.Header("Api-Key", key)

		c.Next()
	}
}

func ApiLanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		language := utils.GetHeader(c, "X-Language", "Api-Language", "Language")
		if language == "" {
			language = "EN"
		}

		c.Set(apiLanguage, language)
		c.Header("X-Language", language)
		c.Header("Api-Language", language)

		c.Next()
	}
}

func ApiDeviceOSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceOS := utils.GetHeader(c, "X-Device-OS", "Api-DeviceOS", "DeviceOS", "Device-OS")

		c.Set(apiDeviceOS, deviceOS)
		c.Header("X-Device-OS", deviceOS)
		c.Header("Api-DeviceOS", deviceOS)

		c.Next()
	}
}

func ApiChannelMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		channel := utils.GetHeader(c, "X-Channel", "Api-Channel", "Channel")

		c.Set(apiChannel, channel)
		c.Header("X-Channel", channel)
		c.Header("Api-Channel", channel)

		c.Next()
	}
}

// PrometheusMiddleware
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		status := fmt.Sprintf("%d", c.Writer.Status())
		path := c.FullPath()
		method := c.Request.Method
		metrics.HttpRequestsTotal.With(prometheus.Labels{"method": method, "path": path, "status": status}).Inc()
		metrics.HttpRequestDuration.With(prometheus.Labels{"method": method, "path": path, "status": status}).Observe(time.Since(start).Seconds())
	}
}

// HealthCheck provides a simple health check endpoint.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
