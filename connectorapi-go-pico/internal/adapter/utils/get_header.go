package utils

import "github.com/gin-gonic/gin"

// GetHeader returns the first non-empty header value found from a list of possible keys
func GetHeader(c *gin.Context, keys ...string) string {
	for _, key := range keys {
		if val := c.GetHeader(key); val != "" {
			return val
		}
	}
	return ""
}