package middleware

import "github.com/Alkemic/go-route"

type Middleware func(route.HandlerFunc) route.HandlerFunc

func Chain(middlewares ...Middleware) Middleware {
	return func(f route.HandlerFunc) route.HandlerFunc {
		for _, middleware := range middlewares {
			f = middleware(f)
		}

		return f
	}
}
