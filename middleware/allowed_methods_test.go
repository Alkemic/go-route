package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alkemic/go-route"
)

func TestAllowedMethods(t *testing.T) {
	basicHandler := func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, "response ok")
		rw.WriteHeader(http.StatusOK)
	}

	testCases := []struct {
		name string

		method         string
		handler        route.HandlerFunc
		allowedMethods []string

		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "allowed get",
			method:             http.MethodGet,
			handler:            basicHandler,
			allowedMethods:     []string{http.MethodGet},
			expectedStatusCode: http.StatusOK,
			expectedBody:       "response ok",
		}, {
			name:               "not allowed get",
			method:             http.MethodGet,
			handler:            basicHandler,
			allowedMethods:     []string{http.MethodPost, http.MethodPut},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedBody:       "405 method not allowed\n",
		}, {
			name:               "not allowed get different case",
			method:             http.MethodGet,
			handler:            basicHandler,
			allowedMethods:     []string{"get"},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedBody:       "405 method not allowed\n",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "http://example.com/foo", nil)
			w := httptest.NewRecorder()

			handler := AllowedMethods(tc.allowedMethods)(tc.handler)

			handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			if tc.expectedStatusCode != resp.StatusCode {
				t.Errorf("Expected status code '%d', but got '%d'", tc.expectedStatusCode, resp.StatusCode)
			}

			if tc.expectedBody != string(body) {
				t.Errorf("Expected body '%s', but got '%s'", tc.expectedBody, string(body))
			}
		})
	}
}
