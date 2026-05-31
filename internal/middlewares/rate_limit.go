package middlewares

import (
	"go-chat/internal/response"
	"net/http"

	"golang.org/x/time/rate"
)                                                                                                                                                                                                       
                                                                                                                                                                                                        
  func RateLimitMiddleware(rps int, burst int) func(http.Handler) http.Handler {
      limiter := rate.NewLimiter(rate.Limit(rps), burst)
      return func(next http.Handler) http.Handler {
          return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {                                                                                                                          
              if !limiter.Allow() {       
                  response.WriteError(w, http.StatusTooManyRequests, "TOO_MANY_REQUESTS", "rate limit exceeded")                                                                                          
                  return                                                                                                                                                                                  
              }                           
              next.ServeHTTP(w, r)                                                                                                                                                                        
          })                                                                                                                                                                                              
      }
  }             