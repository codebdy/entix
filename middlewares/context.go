package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	"rxdrag.com/entify/authentication"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/contexts"
)

// 传递公共参数中间件
func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//为了测试loading状态，生产版需要删掉
		time.Sleep(time.Duration(300) * time.Millisecond)

		reqToken := r.Header.Get(consts.AUTHORIZATION)
		splitToken := strings.Split(reqToken, consts.BEARER)
		v := contexts.ContextValues{}
		if len(splitToken) == 2 {
			reqToken = splitToken[1]
			if reqToken != "" {
				v.Token = reqToken
				me, err := authentication.GetUserByToken(reqToken)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				v.Me = me
			}
		}
		ctx := context.WithValue(r.Context(), consts.CONTEXT_VALUES, v)
		appUuid := r.Header.Get(consts.HEADER_APPX_APPUUID)
		if appUuid == "" {
			appUuid = consts.SYSTEM_APP_UUID
		}
		ctx = context.WithValue(ctx, consts.APPUUID, appUuid)
		ctx = context.WithValue(ctx, consts.HOST, r.Host)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
