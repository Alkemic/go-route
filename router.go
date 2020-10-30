package route

import (
	"net/http"
	"regexp"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	handle(string, http.ResponseWriter, *http.Request)
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (f HandlerFunc) handle(_ string, resp http.ResponseWriter, req *http.Request) {
	f(resp, req)
}

// RegexpRouter
type RegexpRouter struct {
	routes   []route
	NotFound func(w http.ResponseWriter, r *http.Request)
}

func New() *RegexpRouter {
	return &RegexpRouter{
		NotFound: http.NotFound,
	}
}

func (r *RegexpRouter) Add(pattern string, handler interface{}, allowedMethods ...string) {
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
}

func (r RegexpRouter) handle(urlPath string, rw http.ResponseWriter, req *http.Request) {
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
			route.handler.handle(urlPath, rw, req)
			return
		}
	}
	r.NotFound(rw, req)
}

func (r RegexpRouter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	initParams(req)
	r.handle(req.URL.Path, rw, req)
}
