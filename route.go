package route

import (
	"net/http"
	"regexp"
)

var (
	defaultMethods = map[string]struct{}{
		http.MethodGet:     {},
		http.MethodHead:    {},
		http.MethodPost:    {},
		http.MethodPut:     {},
		http.MethodPatch:   {},
		http.MethodDelete:  {},
		http.MethodConnect: {},
		http.MethodOptions: {},
		http.MethodTrace:   {},
	}
)

type route struct {
	pattern        *regexp.Regexp
	handler        Handler
	allowedMethods map[string]struct{}
}

func newRoute(pattern *regexp.Regexp, handler Handler, allowedMethods ...string) route {
	allowedMethodsMap := map[string]struct{}{}
	if len(allowedMethods) == 0 {
		allowedMethodsMap = defaultMethods
	} else {
		for _, method := range allowedMethods {
			allowedMethodsMap[method] = struct{}{}
		}
	}
	return route{
		pattern:        pattern,
		handler:        handler,
		allowedMethods: allowedMethodsMap,
	}
}
