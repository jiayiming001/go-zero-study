package handler

import (
	"net/http"

	"go-zero-study/core/logx"
	"go-zero-study/core/syncx"
	"go-zero-study/rest/internal"
)

func MaxConns(n int) func(http.Handler) http.Handler {
	if n <= 0 {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		latchLimiter := syncx.NewLimit(n)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if latchLimiter.TryBorrow() {
				defer func() {
					if err := latchLimiter.Return(); err != nil {
						logx.Error(err)
					}
				}()

				next.ServeHTTP(w, r)
			} else {
				internal.Errorf(r, "Concurrent connections over %d, rejected with code %d",
					n, http.StatusServiceUnavailable)
				w.WriteHeader(http.StatusServiceUnavailable)
			}
		})
	}
}
