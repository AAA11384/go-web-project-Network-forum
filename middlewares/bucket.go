package middlewares

import (
	"context"
	"qimiproject/controllers"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

var rl ratelimit.Limiter

func BucketLimit(RequestPerSecond int, maxTime time.Duration) func(c *gin.Context) {
	rl = ratelimit.New(RequestPerSecond)
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(context.Background())
		timer := time.NewTimer(maxTime)
		ch := make(chan interface{}, 1)
		go func(ctx context.Context) {
			select {
			case ch <- rl.Take():
				return
			case <-ctx.Done():
				return
			}
		}(ctx)
		select {
		case <-timer.C:
			controllers.ResponseError(c, controllers.CodeRequestFrequently)
			cancel()
			c.Abort()
			return
		case <-ch:
			c.Next()
		}
	}
}
