package middleware

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNoop(t *testing.T) {
	testFunc := func(w http.ResponseWriter, r *http.Request) {}
	output := Noop(testFunc)
	if reflect.ValueOf(testFunc).Pointer() != reflect.ValueOf(output).Pointer() {
		t.Fatalf("Function returned from Noop differs from inputed")
	}
}
