package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alkemic/go-route"
)

func TestChain(t *testing.T) {
	type check func(rec *httptest.ResponseRecorder, t *testing.T)
	checks := func(cs ...check) []check { return cs }

	hasHeader := func(expectedHeaders map[string]string) check {
		return func(rec *httptest.ResponseRecorder, t *testing.T) {
			for k, v := range expectedHeaders {
				h := rec.Result().Header.Get(k)
				if h != v {
					t.Errorf("Expected header '%s' to be '%s', but got '%s'", k, v, h)
				}
			}
		}
	}
	hasBody := func(exp string) check {
		return func(rec *httptest.ResponseRecorder, t *testing.T) {
			if exp != rec.Body.String() {
				t.Errorf("Expected body '%s', but got '%s'", exp, rec.Body.String())
			}
		}
	}
	hasStatus := func(exp int) check {
		return func(rec *httptest.ResponseRecorder, t *testing.T) {
			if rec.Result().StatusCode != exp {
				t.Errorf("Expected satus code '%d', but got '%d'", exp, rec.Result().StatusCode)
			}
		}
	}

	basicHandler := func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("response ok"))
	}

	writeHeader := func(k, v string) route.Middleware {
		return func(f route.HandlerFunc) route.HandlerFunc {
			return func(rw http.ResponseWriter, req *http.Request) {
				rw.Header().Set(k, v)
				f(rw, req)
			}
		}
	}
	writeBodyMiddleware := func(data string) route.Middleware {
		return func(f route.HandlerFunc) route.HandlerFunc {
			return func(rw http.ResponseWriter, req *http.Request) {
				rw.Write([]byte(data))
				f(rw, req)
			}
		}
	}
	setStatusMiddleware := func(code int) route.Middleware {
		return func(f route.HandlerFunc) route.HandlerFunc {
			return func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(code)
				f(rw, req)
			}
		}
	}
	breakChainMiddleware := func(f route.HandlerFunc) route.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusMovedPermanently)
		}
	}

	testCases := []struct {
		name        string
		handler     route.HandlerFunc
		middlewares []route.Middleware
		checks      []check
	}{
		{
			name:    "no middlewares",
			handler: basicHandler,
			checks: checks(
				hasStatus(http.StatusOK),
				hasBody("response ok"),
			),
		}, {
			name:        "single middleware",
			handler:     basicHandler,
			middlewares: []route.Middleware{writeHeader("Header-Middleware1", "value1")},
			checks: checks(
				hasStatus(http.StatusOK),
				hasBody("response ok"),
				hasHeader(map[string]string{"Header-Middleware1": "value1"}),
			),
		}, {
			name:        "multiple middleware",
			handler:     basicHandler,
			middlewares: []route.Middleware{writeHeader("Header-Middleware1", "value1"), writeHeader("Header-Middleware2", "value2")},
			checks: checks(
				hasStatus(http.StatusOK),
				hasBody("response ok"),
				hasHeader(map[string]string{"Header-Middleware1": "value1", "Header-Middleware2": "value2"}),
			),
		}, {
			name:        "multiple middleware reverse order",
			handler:     basicHandler,
			middlewares: []route.Middleware{writeHeader("Header-Middleware2", "value2"), writeHeader("Header-Middleware1", "value1")},
			checks: checks(
				hasStatus(http.StatusOK),
				hasBody("response ok"),
				hasHeader(map[string]string{"Header-Middleware1": "value1", "Header-Middleware2": "value2"}),
			),
		}, {
			name:        "break chain",
			handler:     basicHandler,
			middlewares: []route.Middleware{writeHeader("Header-Middleware2", "value2"), breakChainMiddleware, writeHeader("Header-Middleware1", "value1")},
			checks: checks(
				hasStatus(http.StatusMovedPermanently),
				hasBody(""),
				hasHeader(map[string]string{"Header-Middleware2": "value2"}),
			),
		}, {
			name:        "writing body test",
			handler:     basicHandler,
			middlewares: []route.Middleware{writeBodyMiddleware("spam\n"), writeBodyMiddleware("ham\n")},
			checks: checks(
				hasStatus(http.StatusOK),
				hasBody("spam\nham\nresponse ok"),
			),
		}, {
			name:        "writing body test in reverse order",
			handler:     basicHandler,
			middlewares: []route.Middleware{writeBodyMiddleware("ham\n"), writeBodyMiddleware("spam\n")},
			checks: checks(
				hasStatus(http.StatusOK),
				hasBody("ham\nspam\nresponse ok"),
			),
		}, {
			name:        "writing status code from middleware",
			handler:     basicHandler,
			middlewares: []route.Middleware{setStatusMiddleware(http.StatusAccepted)},
			checks: checks(
				hasStatus(http.StatusAccepted),
				hasBody("response ok"),
			),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			rec := httptest.NewRecorder()
			Chain(tc.middlewares...)(tc.handler)(rec, req)

			for _, ch := range tc.checks {
				ch(rec, t)
			}
		})
	}
}
