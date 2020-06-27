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

type route struct {
	pattern *regexp.Regexp
	handler Handler
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

func (h *RegexpRouter) Add(pattern string, handler interface{}) {
	var handlerFunc Handler

	switch _handler := handler.(type) {
	case func(http.ResponseWriter, *http.Request):
		handlerFunc = HandlerFunc(_handler)
	case http.HandlerFunc:
	case HandlerFunc:
	case RegexpRouter:
	case *RegexpRouter:
		handlerFunc = _handler
	default:
		panic("Unknown handler param passed to RegexpRouter.Add")
	}

	h.routes = append(h.routes, route{pattern: regexp.MustCompile(pattern), handler: handlerFunc})
}

func (h RegexpRouter) handle(urlPath string, resp http.ResponseWriter, req *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(urlPath) {
			match := route.pattern.FindStringSubmatch(urlPath)
			for i, name := range route.pattern.SubexpNames() {
				if i != 0 {
					SetParam(req, name, match[i])
				}
			}
			urlPath = route.pattern.ReplaceAllString(urlPath, "")
			route.handler.handle(urlPath, resp, req)
			return
		}
	}
	h.NotFound(resp, req)
}

func (h RegexpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	initParams(r)
	h.handle(r.URL.Path, w, r)
}
