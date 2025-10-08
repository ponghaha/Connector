package client
import (
	"io"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	//"connectorapi-go/internal/core/domain"

	appError "connectorapi-go/pkg/error"
)
type HTTPClient struct {
	client *http.Client
	logger *zap.SugaredLogger
}
func NewHTTPClient(logger *zap.SugaredLogger) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.logger.Infow("Sending downstream HTTP request",
		"method", req.Method,
		"url", req.URL.String(),
	)
	return c.client.Do(req)
}
func (c *HTTPClient) Forward(ctx *gin.Context, targetURL, apiKey string) *appError.AppError {
	c.logger.Infow("Forwarding HTTP request", "targetURL", targetURL)
	req, err := http.NewRequestWithContext(ctx.Request.Context(), ctx.Request.Method, targetURL, ctx.Request.Body)
	if err != nil {
		c.logger.Errorw("Failed to create proxy request", "error", err)
		return appError.ErrInternalServer
	}
	req.Header = ctx.Request.Header.Clone()
	if apiKey != "" {
		req.Header.Set("API-Key", apiKey)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Errorw("Failed to call downstream service for proxy", "error", err)
		return appError.ErrService
	}
	defer resp.Body.Close()
	for key, values := range resp.Header {
		for _, value := range values {
			ctx.Writer.Header().Add(key, value)
		}
	}
	ctx.Writer.WriteHeader(resp.StatusCode)
	io.Copy(ctx.Writer, resp.Body)
	return nil
}