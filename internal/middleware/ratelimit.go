package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

type window struct {
	count   int
	resetAt time.Time
}

// RateLimit membatasi jumlah request per key dalam satu window waktu (in-memory, fixed window).
func RateLimit(limit int, per time.Duration, key func(*gin.Context) string) gin.HandlerFunc {
	var mu sync.Mutex
	windows := make(map[string]*window)

	return func(c *gin.Context) {
		k := key(c)
		now := time.Now()

		mu.Lock()
		w, ok := windows[k]
		if !ok || now.After(w.resetAt) {
			w = &window{resetAt: now.Add(per)}
			windows[k] = w
			if len(windows) > 10000 {
				for wk, v := range windows {
					if now.After(v.resetAt) {
						delete(windows, wk)
					}
				}
			}
		}
		w.count++
		exceeded := w.count > limit
		retryAfter := int(w.resetAt.Sub(now).Seconds()) + 1
		mu.Unlock()

		if exceeded {
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			utils.Error(c, http.StatusTooManyRequests, "too many requests, please try again later")
			c.Abort()
			return
		}
		c.Next()
	}
}

func ByIP(c *gin.Context) string { return "ip:" + c.ClientIP() }

func ByUser(c *gin.Context) string { return fmt.Sprintf("user:%d", CurrentUserID(c)) }
