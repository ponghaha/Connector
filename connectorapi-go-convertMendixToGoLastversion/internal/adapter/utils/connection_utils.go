package utils

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	rng   = rand.New(rand.NewSource(time.Now().UnixNano()))
	rngMu sync.Mutex
)

// GetRouteKey returns the route key in the format METHOD:/path
func GetRouteKey(c *gin.Context) string {
	return c.Request.Method + ":" + c.FullPath()
	
}

// RandomPortFromList returns a random port from a list of ports.
// Returns empty string if the list is empty.
func RandomPortFromList(portList []string) string {
	if len(portList) == 0 {
		return ""
	}
	rngMu.Lock()
	defer rngMu.Unlock()
	return portList[rng.Intn(len(portList))]
}
