package middleware

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/Alkemic/go-route"
)

var (
	UserKey = "_user_id_"

	ErrUnknownUser           = errors.New("unknown user")
	ErrEmptyList             = errors.New("empty user list provided")
	ErrNoUserProvided        = errors.New("no user provided in request context")
	ErrAuthHeaderWrongLength = errors.New("wrong length of authorization value")
)

type AuthFn func(string, string) (string, error)

// BasicAuthenticate setups authenticate headers in request, and handles authenticate based on authFn. In case
// of valid user authentication, it also set user id (that was returned from authFn) in request's context, and
// can be retrieved by GetUser function.
func BasicAuthenticate(logger *log.Logger, authFn AuthFn, realm string) func(http.HandlerFunc) http.HandlerFunc {
	return func(view http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, password, err := getCredentials(r)
			if err != nil {
				logger.Println(err)
			}
			userID, err := authFn(user, password)
			if err != nil {
				w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				return
			}

			route.SetParam(r, UserKey, userID)

			view(w, r)
		}
	}
}

func getCredentials(r *http.Request) (string, string, error) {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return "", "", errors.New("missing authorization header")
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return "", "", errors.Wrap(err, "can't decode authorization header")
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return "", "", ErrAuthHeaderWrongLength
	}

	return pair[0], pair[1], nil
}

// Authenticate verify credentials against provided user and password.
func Authenticate(expectedUser, expectedPassword string) AuthFn {
	return func(user, password string) (string, error) {
		if user == expectedUser && expectedPassword == password {
			return user, nil
		}

		return "", ErrUnknownUser
	}
}

// Authenticate verify credentials against provided map (user => password).
func AuthenticateMap(userList map[string]string) AuthFn {
	return func(user, password string) (string, error) {
		if len(userList) == 0 {
			return "", ErrEmptyList
		}
		for u, p := range userList {
			if u == user && p == password {
				return user, nil
			}
		}
		return "", ErrUnknownUser
	}
}

// GetUser returns user from request's context, and error when user is not provided.
func GetUser(r *http.Request) (string, error) {
	user, ok := route.GetParam(r, UserKey)
	if !ok {
		return "", ErrNoUserProvided
	}
	return user, nil
}
