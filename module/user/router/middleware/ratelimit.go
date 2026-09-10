package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	commonresponse "fastduck/treasure-doc/module/common/response"

	"github.com/gin-gonic/gin"
)

// RateRule 描述某个路径前缀的限流参数：Rate 为每秒补充的令牌数，Burst 为桶容量。
type RateRule struct {
	Prefix string
	Rate   float64
	Burst  float64
}

type tokenBucket struct {
	tokens float64
	last   time.Time
}

// bucketTTL 内无访问的桶会被清理，避免内存随 IP 无限增长。
const bucketTTL = 10 * time.Minute

// RateLimit 按客户端 IP 做令牌桶限流：命中前缀的请求才受限，其余直接放行。
func RateLimit(rules []RateRule) gin.HandlerFunc {
	var (
		mu      sync.Mutex
		buckets = make(map[string]*tokenBucket)
	)

	go func() {
		ticker := time.NewTicker(bucketTTL)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			for key, bucket := range buckets {
				if time.Since(bucket.last) > bucketTTL {
					delete(buckets, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		rule := matchRule(rules, c.Request.URL.Path)
		if rule == nil {
			c.Next()
			return
		}

		key := rule.Prefix + "|" + c.ClientIP()
		now := time.Now()

		mu.Lock()
		bucket, ok := buckets[key]
		if !ok {
			bucket = &tokenBucket{tokens: rule.Burst, last: now}
			buckets[key] = bucket
		}
		bucket.tokens = min(rule.Burst, bucket.tokens+now.Sub(bucket.last).Seconds()*rule.Rate)
		bucket.last = now
		allowed := bucket.tokens >= 1
		if allowed {
			bucket.tokens--
		}
		mu.Unlock()

		if !allowed {
			commonresponse.Error(c, http.StatusTooManyRequests, 42900, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}

func matchRule(rules []RateRule, path string) *RateRule {
	for index := range rules {
		if strings.HasPrefix(path, rules[index].Prefix) {
			return &rules[index]
		}
	}
	return nil
}
