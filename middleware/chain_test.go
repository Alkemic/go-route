package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alkemic/go-route"
)

func TestChain(t *testing.T) {
	basicHandler := func(rw http.ResponseWriter, req *http.Request, p map[string]string) {
		rw.Write([]byte("response ok"))
		rw.WriteHeader(http.StatusOK)
	}

	middleware1 := func(f route.HandlerFunc) route.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request, params map[string]string) {
			rw.Header().Set("Header-Middleware1", "value1")
		}
	}
	middleware2 := func(f route.HandlerFunc) route.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request, params map[string]string) {
			rw.Header().Set("Header-Middleware2", "value2")
		}
	}

	testCases := []struct {
		name            string
		handler         route.HandlerFunc
		middlewares     []Middleware
		expectedBody    string
		expectedHeaders map[string]string
	}{
		{
			name:         "no middlewares",
			handler:      basicHandler,
			expectedBody: "response ok",
		}, {
			name:            "single middleware",
			handler:         basicHandler,
			middlewares:     []Middleware{middleware1},
			expectedBody:    "response ok",
			expectedHeaders: map[string]string{"Header-Middleware1": "value1"},
		}, {
			name:            "middleware order",
			handler:         basicHandler,
			middlewares:     []Middleware{middleware1, middleware2},
			expectedBody:    "response ok",
			expectedHeaders: map[string]string{"Header-Middleware1": "value1", "Header-Middleware2": "value2"},
		}, {
			name:            "middleware order",
			handler:         basicHandler,
			middlewares:     []Middleware{middleware2, middleware1},
			expectedBody:    "response ok",
			expectedHeaders: map[string]string{"Header-Middleware1": "value1", "Header-Middleware2": "value2"},
		}, {
			name:            "single middleware",
			handler:         middleware1(basicHandler),
			middlewares:     nil,
			expectedBody:    "response ok",
			expectedHeaders: map[string]string{"Header-Middleware1": "value1"},
		}, {
			name:            "middleware order",
			handler:         middleware1(middleware2(basicHandler)),
			middlewares:     nil,
			expectedBody:    "response ok",
			expectedHeaders: map[string]string{"Header-Middleware1": "value1", "Header-Middleware2": "value2"},
		}, {
			name:            "middleware order",
			handler:         middleware2(middleware1(basicHandler)),
			middlewares:     nil,
			expectedBody:    "response ok",
			expectedHeaders: map[string]string{"Header-Middleware1": "value1", "Header-Middleware2": "value2"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			w := httptest.NewRecorder()

			Chain(tc.middlewares...)(tc.handler)(w, req, map[string]string{})

			resp := w.Result()
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			if tc.expectedBody != string(body) {
				t.Errorf("Expected body '%s', but got '%s'", tc.expectedBody, string(body))
			}

			for k, v := range tc.expectedHeaders {
				h := resp.Header.Get(k)
				if h != v {
					t.Errorf("Expected header '%s' to be '%s', but got '%s'", k, v, h)
				}
			}

		})
	}
}
