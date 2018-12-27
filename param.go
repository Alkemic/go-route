package route

import (
	"context"
	"net/http"
)

var paramsKey = "_params_"

// GetParams returns map of parameters parsed from url, will be empty in case if there were no params.
func GetParams(r *http.Request) map[string]string {
	params, _ := r.Context().Value(paramsKey).(map[string]string)
	return params
}

func SetParam(r *http.Request, key, value string) {
	ctx := r.Context()
	params, _ := ctx.Value(paramsKey).(map[string]string)
	if params == nil {
		initParams(r)
		ctx = r.Context()
		params, _ = ctx.Value(paramsKey).(map[string]string)
	}
	params[key] = value
	ctx = context.WithValue(ctx, paramsKey, params)
	(*r) = *(r.WithContext(ctx))
}

func initParams(r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, paramsKey, map[string]string{})
	(*r) = *(r.WithContext(ctx))
}

// GetParam returns single parsed parameter from url, along with indicator if it was success.
func GetParam(r *http.Request, key string) (string, bool) {
	params := GetParams(r)
	val, ok := params[key]
	return val, ok
}
