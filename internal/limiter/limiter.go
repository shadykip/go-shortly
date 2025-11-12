package limiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type IPLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.RWMutex
}

func NewIPLimiter(r rate.Limit, b int) *IPLimiter {
	return &IPLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  sync.RWMutex{},
	}
}

func (i *IPLimiter) getLimiter(ip string, r rate.Limit, b int) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(r, b)
		i.ips[ip] = limiter
	}
	return limiter
}

// Allow checks if request is allowed
func (i *IPLimiter) Allow(ip string, r rate.Limit, b int) bool {
	limiter := i.getLimiter(ip, r, b)
	return limiter.Allow()
}

// RetryAfter returns seconds to wait (if denied)
func (i *IPLimiter) RetryAfter(ip string, r rate.Limit, b int) time.Duration {
	limiter := i.getLimiter(ip, r, b)
	return limiter.Reserve().Delay()
}
