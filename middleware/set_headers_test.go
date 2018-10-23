package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alkemic/go-route"
)

func TestSetHeaders(t *testing.T) {
	basicHandler := func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, "response ok")
		rw.WriteHeader(http.StatusOK)
	}

	testCases := []struct {
		name string

		handler route.HandlerFunc
		headers map[string]string

		expectedStatusCode int
		expectedHeaders    map[string]string
	}{
		{
			name:               "set basic headers",
			handler:            basicHandler,
			headers:            map[string]string{"Header1": "Value1"},
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    map[string]string{"Header1": "Value1", "Content-Type": "text/plain; charset=utf-8"},
		}, {
			name:               "override content type",
			handler:            basicHandler,
			headers:            map[string]string{"Content-Type": "application/json"},
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    map[string]string{"Content-Type": "application/json"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			w := httptest.NewRecorder()

			handler := SetHeaders(tc.headers)(tc.handler)

			handler(w, req)

			resp := w.Result()

			if tc.expectedStatusCode != resp.StatusCode {
				t.Errorf("Expected status code '%d', but got '%d'", tc.expectedStatusCode, resp.StatusCode)
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
