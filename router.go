package route

import (
	"context"
	"net/http"
	"regexp"
)

var urlPathContextKey = struct{}{}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	handle(http.ResponseWriter, *http.Request)
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	f(rw, req)
}

func (f HandlerFunc) handle(rw http.ResponseWriter, req *http.Request) {
	f(rw, req)
}

type Middleware func(fn HandlerFunc) HandlerFunc

// RegexpRouter
type RegexpRouter struct {
	routes      []route
	middlewares []Middleware
	NotFound    func(w http.ResponseWriter, r *http.Request)
}

func New() *RegexpRouter {
	return &RegexpRouter{
		NotFound: http.NotFound,
	}
}

func (r *RegexpRouter) Add(pattern string, handler interface{}, allowedMethods ...string) *RegexpRouter {
	var handlerFunc Handler

	switch _handler := handler.(type) {
	case func(http.ResponseWriter, *http.Request):
		handlerFunc = HandlerFunc(_handler)
	case http.HandlerFunc:
		handlerFunc = HandlerFunc(_handler)
	case RegexpRouter:
		handlerFunc = _handler
	case *RegexpRouter:
		handlerFunc = _handler
	default:
		panic("Unknown handler param passed to RegexpRouter.Add")
	}

	r.routes = append(r.routes, newRoute(regexp.MustCompile(pattern), handlerFunc, allowedMethods...))

	return r
}

func (r *RegexpRouter) AddMiddleware(mw Middleware) *RegexpRouter {
	r.middlewares = append(r.middlewares, mw)
	return r
}

func (r RegexpRouter) handle(rw http.ResponseWriter, req *http.Request) {
	urlPath := req.Context().Value(urlPathContextKey).(string)
	for _, route := range r.routes {
		if route.pattern.MatchString(urlPath) {
			if _, ok := route.allowedMethods[req.Method]; !ok {
				http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			match := route.pattern.FindStringSubmatch(urlPath)
			for i, name := range route.pattern.SubexpNames() {
				if i != 0 {
					SetParam(req, name, match[i])
				}
			}
			urlPath = route.pattern.ReplaceAllString(urlPath, "")
			req = req.WithContext(context.WithValue(req.Context(), urlPathContextKey, urlPath))
			fn := route.handler.handle
			for _, middleware := range r.middlewares {
				fn = middleware(fn)
			}
			fn(rw, req)
			return
		}
	}
	r.NotFound(rw, req)
}

func (r RegexpRouter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	initParams(req)
	req = req.WithContext(context.WithValue(req.Context(), urlPathContextKey, req.URL.Path))
	r.handle(rw, req)
}
