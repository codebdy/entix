package app

import (
	"context"
	"net/http"

	"rxdrag.com/entify/app/resolve"
	"rxdrag.com/entify/consts"
)

func LoadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), consts.LOADERS, resolve.CreateDataLoaders())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
