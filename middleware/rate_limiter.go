package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func PerClientTokenBucket() gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter // that applies token bucket
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client) // map[IP]client
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return gin.HandlerFunc(func(ctx *gin.Context) {
		// IP addess of request client
		ip, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		mu.Lock()
		if _, found := clients[ip]; !found {
			// add new client
			// 1 requests/sec, burst size = 3
			clients[ip] = &client{limiter: rate.NewLimiter(1, 3)}
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests from the client"})
			ctx.Abort() // Abort the request pipeline
			// otherwise the request will still be processed
			return
		} else {
			mu.Unlock()
			ctx.Next()
		}
	})
}
