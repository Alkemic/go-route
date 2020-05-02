package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	basicHandler = func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, "response ok")
		rw.WriteHeader(http.StatusOK)
	}
	stringPanicHandler = func(rw http.ResponseWriter, req *http.Request) {
		panic("string error")
	}
	errorPanicHandler = func(rw http.ResponseWriter, req *http.Request) {
		panic(errors.New("error error"))
	}

	testCases = []struct {
		name               string
		handler            http.HandlerFunc
		expectedStatusCode int
		expectedBody       string
		expectedLog        string
	}{
		{
			name:               "handler without panic",
			handler:            basicHandler,
			expectedStatusCode: http.StatusOK,
			expectedBody:       "response ok",
		}, {
			name:               "panic with string",
			handler:            stringPanicHandler,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "500 internal server error\n",
			expectedLog:        "Panic occured:\nstring error\nPanic end.\n",
		}, {
			name:               "panic with error",
			handler:            errorPanicHandler,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "500 internal server error\n",
			expectedLog:        "Panic occured:\nerror error\nPanic end.\n",
		},
	}
)

func TestPanicInterceptor(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			writer := io.Writer(&buffer)
			log.SetOutput(writer)
			flags := log.Flags()
			log.SetFlags(0)
			defer func() {
				log.SetOutput(os.Stderr)
				log.SetFlags(flags)
			}()

			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			w := httptest.NewRecorder()

			PanicInterceptor(tc.handler)(w, req)

			resp := w.Result()
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			if tc.expectedStatusCode != resp.StatusCode {
				t.Errorf("Expected status code '%d', but got '%d'", tc.expectedStatusCode, resp.StatusCode)
			}

			if tc.expectedLog != buffer.String() {
				t.Errorf("Expected log output '%s', but got '%s'", tc.expectedLog, buffer.String())
			}

			if tc.expectedBody != string(body) {
				t.Errorf("Expected body '%s', but got '%s'", tc.expectedBody, string(body))
			}
		})
	}
}

func TestPanicInterceptorWithLogger(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			logger := log.New(io.Writer(&buffer), "", 0)

			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			w := httptest.NewRecorder()

			PanicInterceptorWithLogger(logger)(tc.handler)(w, req)

			resp := w.Result()
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			if tc.expectedStatusCode != resp.StatusCode {
				t.Errorf("Expected status code '%d', but got '%d'", tc.expectedStatusCode, resp.StatusCode)
			}

			if tc.expectedLog != buffer.String() {
				t.Errorf("Expected log output '%s', but got '%s'", tc.expectedLog, buffer.String())
			}

			if tc.expectedBody != string(body) {
				t.Errorf("Expected body '%s', but got '%s'", tc.expectedBody, string(body))
			}
		})
	}
}
