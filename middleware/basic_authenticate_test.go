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
	"net/textproto"
	"testing"
)

func TestAuthenticateList(t *testing.T) {
	type check func(string, error, *testing.T)
	checks := func(cs ...check) []check { return cs }

	hasError := func(msg string) check {
		return func(_ string, err error, t *testing.T) {
			t.Helper()
			if err == nil {
				t.Fatalf("Expected error: '%s' but got nil.", msg)
			}
			if err.Error() != msg {
				t.Errorf("Expected error '%v', got: '%s'", msg, err.Error())
			}
		}
	}
	hasNoError := func(_ string, err error, t *testing.T) {
		t.Helper()
		if err != nil {
			t.Errorf("Unexpected error: '%v'", err)
		}
	}
	hasUserID := func(exp string) check {
		return func(user string, err error, t *testing.T) {
			t.Helper()
			if exp != user {
				t.Errorf("Expected user to be '%s', got: '%s'", user, exp)
			}
		}
	}

	testCases := []struct {
		name     string
		users    map[string]string
		user     string
		password string
		checks   []check
	}{
		{
			name:     "fail on empty list",
			user:     "john",
			password: "asdadasdasd",
			checks: checks(
				hasUserID(""),
				hasError("empty user list provided"),
			),
		}, {
			name: "success on list with one correct user",
			users: map[string]string{
				"john": "asdadasdasd",
			},
			user:     "john",
			password: "asdadasdasd",
			checks: checks(
				hasUserID("john"),
				hasNoError,
			),
		}, {
			name: "fail on list with one incorrect user",
			users: map[string]string{
				"john": "asdadasdasd1",
			},
			user:     "john",
			password: "asdadasdasd",
			checks: checks(
				hasUserID(""),
				hasError("unknown user"),
			),
		}, {
			name: "fail on list with one incorrect user",
			users: map[string]string{
				"john1": "asdadasdasd",
			},
			user:     "john",
			password: "asdadasdasd",
			checks: checks(
				hasUserID(""),
				hasError("unknown user"),
			),
		}, {
			name: "fail on list with one incorrect user",
			users: map[string]string{
				"john1": "asdadasdasd",
				"john":  "asdadasdasd",
				"john3": "asdadasdasd",
				"john2": "asdadasdasd",
			},
			user:     "john",
			password: "asdadasdasd",
			checks: checks(
				hasUserID("john"),
				hasNoError,
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userID, err := AuthenticateMap(tc.users)(tc.user, tc.password)
			for _, ch := range tc.checks {
				ch(userID, err, t)
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	type check func(string, error, *testing.T)
	checks := func(cs ...check) []check { return cs }

	hasError := func(msg string) check {
		return func(_ string, err error, t *testing.T) {
			t.Helper()
			if err == nil {
				t.Fatalf("Expected error: '%s' but got nil.", msg)
			}
			if err.Error() != msg {
				t.Errorf("Expected error '%v', got: '%s'", msg, err.Error())
			}
		}
	}
	hasNoError := func(_ string, err error, t *testing.T) {
		t.Helper()
		if err != nil {
			t.Errorf("Unexpected error: '%v'", err)
		}
	}
	hasUserID := func(exp string) check {
		return func(user string, err error, t *testing.T) {
			t.Helper()
			if exp != user {
				t.Errorf("Expected user to be '%s', got: '%s'", user, exp)
			}
		}
	}

	testCases := []struct {
		name             string
		expectedUser     string
		expectedPassword string
		user             string
		password         string
		checks           []check
	}{
		{
			name:             "success on list with one correct user",
			expectedUser:     "john",
			expectedPassword: "asdadasdasd",
			user:             "john",
			password:         "asdadasdasd",
			checks: checks(
				hasUserID("john"),
				hasNoError,
			),
		}, {
			name:             "fail on incorrect password",
			expectedUser:     "john",
			expectedPassword: "asdadasdasd1",
			user:             "john",
			password:         "asdadasdasd",
			checks: checks(
				hasUserID(""),
				hasError("unknown user"),
			),
		}, {
			name:             "fail on incorrect user",
			expectedUser:     "john",
			expectedPassword: "asdadasdasd",
			user:             "john1",
			password:         "asdadasdasd",
			checks: checks(
				hasUserID(""),
				hasError("unknown user"),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userID, err := Authenticate(tc.expectedUser, tc.expectedPassword)(tc.user, tc.password)
			for _, ch := range tc.checks {
				ch(userID, err, t)
			}
		})
	}
}

func TestBasicAuthenticate(t *testing.T) {
	type check func(rr *httptest.ResponseRecorder, r *http.Request, loggerBuffer bytes.Buffer, t *testing.T)
	checks := func(cs ...check) []check { return cs }

	basicHandler := func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, "response ok")
		rw.WriteHeader(http.StatusOK)
	}

	hasStatusCode := func(code int) check {
		return func(rw *httptest.ResponseRecorder, r *http.Request, loggerBuffer bytes.Buffer, t *testing.T) {
			t.Helper()
			if rw.Code != code {
				t.Fatalf("Expecting status code to be '%d', got '%d'", code, rw.Code)
			}
		}
	}
	hasHeader := func(key, value string) check {
		key = textproto.CanonicalMIMEHeaderKey(key)
		return func(rw *httptest.ResponseRecorder, r *http.Request, loggerBuffer bytes.Buffer, t *testing.T) {
			t.Helper()
			v, ok := rw.HeaderMap[key]
			if !ok {
				t.Fatalf("Key '%s' not found in headers", key)
			}
			if len(v) != 1 {
				t.Fatalf("Value for header '%s', got more than one values: '%v'", key, v)
			}
			if v[0] != value {
				t.Fatalf("Hedaer missmatch, expecting '%s', got '%s'", value, v[0])
			}
		}
	}
	hasNoHeader := func(key string) check {
		key = textproto.CanonicalMIMEHeaderKey(key)
		return func(rw *httptest.ResponseRecorder, r *http.Request, loggerBuffer bytes.Buffer, t *testing.T) {
			t.Helper()
			if _, ok := rw.HeaderMap[key]; ok {
				t.Fatalf("Key '%s' shouldn;t be in headers", key)
			}
		}
	}
	hasBuffer := func(expected string) check {
		return func(rw *httptest.ResponseRecorder, r *http.Request, loggerBuffer bytes.Buffer, t *testing.T) {
			t.Helper()
			if loggerBuffer.String() != expected {
				t.Fatalf("Expecting buffer to be '%s', but got '%s'", expected, loggerBuffer.String())
			}
		}
	}
	hasUser := func(expectedUser string, expectedError error) check {
		return func(rw *httptest.ResponseRecorder, r *http.Request, loggerBuffer bytes.Buffer, t *testing.T) {
			t.Helper()
			user, err := GetUser(r)
			if user != expectedUser {
				t.Fatalf("Expecting user to be '%s', but got '%s'", expectedUser, user)
			}
			if err != expectedError {
				t.Errorf("Expected error to bet '%s', got: %v", expectedError, err)
			}
		}
	}
	hasBody := func(expectedBody string) check {
		return func(rw *httptest.ResponseRecorder, r *http.Request, loggerBuffer bytes.Buffer, t *testing.T) {
			t.Helper()
			resp := rw.Result()
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			if string(body) != expectedBody {
				t.Fatalf("Expecting body to be '%s', but got '%s'", expectedBody, body)
			}
		}
	}

	testAuthenticateFn := func(u string, p string) (string, error) {
		if u == "user" && p == "pass" {
			return "user", nil
		}
		return "", errors.New("mocked error")
	}

	testCases := []struct {
		name           string
		handler        http.HandlerFunc
		headers        http.Header
		authenticateFn AuthFn
		checks         []check
	}{
		{
			name:           "401 when no authenticate header",
			handler:        basicHandler,
			authenticateFn: testAuthenticateFn,
			checks: checks(
				hasStatusCode(401),
				hasHeader("WWW-Authenticate", `Basic realm="test realm"`),
				hasBuffer("missing authorization header\n"),
				hasUser("", ErrNoUserProvided),
				hasBody("401 Unauthorized\n"),
			),
		}, {
			name:           "401 when invalid base 64 header value",
			handler:        basicHandler,
			authenticateFn: testAuthenticateFn,
			headers:        http.Header{"Authorization": []string{"Basic wrong header"}},
			checks: checks(
				hasStatusCode(401),
				hasHeader("WWW-Authenticate", `Basic realm="test realm"`),
				hasBuffer("can't decode authorization header: illegal base64 data at input byte 5\n"),
				hasUser("", ErrNoUserProvided),
				hasBody("401 Unauthorized\n"),
			),
		}, {
			name:           "401 when invalid header value",
			handler:        basicHandler,
			authenticateFn: testAuthenticateFn,
			headers:        http.Header{"Authorization": []string{"Basic dXNlcg=="}},
			checks: checks(
				hasStatusCode(401),
				hasHeader("WWW-Authenticate", `Basic realm="test realm"`),
				hasBuffer("wrong length of authorization value\n"),
				hasUser("", ErrNoUserProvided),
				hasBody("401 Unauthorized\n"),
			),
		}, {
			name:           "401 when no valid user",
			handler:        basicHandler,
			authenticateFn: testAuthenticateFn,
			headers:        http.Header{"Authorization": []string{"Basic dXNlcjphc2RzYWQ="}}, // user:asdasd
			checks: checks(
				hasStatusCode(401),
				hasHeader("WWW-Authenticate", `Basic realm="test realm"`),
				hasUser("", ErrNoUserProvided),
				hasBody("401 Unauthorized\n"),
			),
		}, {
			name:           "ok when valid password",
			handler:        basicHandler,
			authenticateFn: testAuthenticateFn,
			headers:        http.Header{"Authorization": []string{"Basic dXNlcjpwYXNz"}}, // user:pass
			checks: checks(
				hasStatusCode(200),
				hasNoHeader("WWW-Authenticate"),
				hasUser("user", nil),
				hasBody("response ok"),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			logger := log.New(io.Writer(&buffer), "", 0)

			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			req.Header = tc.headers
			w := httptest.NewRecorder()

			handler := BasicAuthenticate(logger, tc.authenticateFn, "test realm")(tc.handler)

			handler(w, req)

			for _, ch := range tc.checks {
				ch(w, req, buffer, t)
			}
		})
	}
}
