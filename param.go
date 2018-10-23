package route

import (
	"context"
	"net/http"
)

var paramsKey = "_params_"

func GetParams(r *http.Request) map[string]string {
	params, _ := r.Context().Value(paramsKey).(map[string]string)
	return params
}

func addParam(r *http.Request, key, value string) {
	ctx := r.Context()
	params, _ := ctx.Value(paramsKey).(map[string]string)
	params[key] = value
	ctx = context.WithValue(ctx, paramsKey, params)
	(*r) = *(r.WithContext(ctx))
}

func initParams(r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, paramsKey, map[string]string{})
	(*r) = *(r.WithContext(ctx))
}
