package middleware

import (
	"github.com/Alkemic/go-route"
)

func Chain(middlewares ...route.Middleware) route.Middleware {
	return func(f route.HandlerFunc) route.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			f = middlewares[i](f)
		}
		return f
	}
}
