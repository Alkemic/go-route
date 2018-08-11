package route

import (
	"net/http"
	"regexp"
)

var (
	NotFound            = http.NotFound
	InternalServerError = func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	}
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	handle(string, http.ResponseWriter, *http.Request, *map[string]string)
}

// Acctual function
type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (f HandlerFunc) handle(
	_ string,
	resp http.ResponseWriter,
	req *http.Request,
	params *map[string]string,
) {
	f(resp, req, *params)
}

type route struct {
	pattern *regexp.Regexp
	handler Handler
}

// RegexpRouter
type RegexpRouter struct {
	routes []*route
}

func (h *RegexpRouter) Add(
	pattern string,
	handler interface{},
) {
	var r *route

	switch _handler := handler.(type) {
	case func(http.ResponseWriter, *http.Request, map[string]string):
		r = &route{regexp.MustCompile(pattern), HandlerFunc(_handler)}
	case HandlerFunc:
	case RegexpRouter:
		r = &route{regexp.MustCompile(pattern), _handler}
	default:
		panic("Unknown handler param passed to RegexpRouter.Add")
	}

	h.routes = append(h.routes, r)
}

func (h RegexpRouter) handle(
	urlPath string,
	resp http.ResponseWriter,
	req *http.Request,
	params *map[string]string,
) {
	for _, route := range h.routes {
		if route.pattern.MatchString(urlPath) {
			match := route.pattern.FindStringSubmatch(urlPath)
			for i, name := range route.pattern.SubexpNames() {
				if i != 0 {
					(*params)[name] = match[i]
				}
			}
			urlPath = route.pattern.ReplaceAllString(urlPath, "")
			route.handler.handle(urlPath, resp, req, params)
			return
		}
	}

	NotFound(resp, req)
}

func (h RegexpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handle(r.URL.Path, w, r, &map[string]string{})
}
