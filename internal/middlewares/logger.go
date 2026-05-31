package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

  func LoggerMiddleware(next http.Handler) http.Handler {                                                                                                                                                                  
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          start := time.Now()                                                                                                                                                                                              
          ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)  
          requestID, _ := r.Context().Value(RequestIDHeader).(string)

          next.ServeHTTP(ww, r)                                                                                                                                                                                            
          slog.Info("request",                                                                                                                                                                                           
              "method", r.Method,                                                                                                                                                                                          
              "path", r.URL.Path,                                                                                                                                                                                          
              "status", ww.Status(),          
              "duration", time.Since(start),
			  "ip", r.RemoteAddr,
              "request_id", requestID,
          )                                                                                                                                                                                                              
      })                                                                                                                                                                                                                   
  }
