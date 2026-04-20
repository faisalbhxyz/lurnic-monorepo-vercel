// Package handler implements the Vercel Go serverless entrypoint for the Gin API.
package handler

import (
	"net/http"
	"strings"

	"dashlearn/pkg/server"
)

const handlerVersion = "v1.0.24"

var ginHandler http.Handler

func init() {
	eng, _, err := server.NewEngine(handlerVersion)
	if err != nil {
		panic(err)
	}
	ginHandler = eng
}

// Handler is invoked by Vercel for each request to /api/vercel/gin-handler.
// Rewrites map /v1/* -> /api/vercel/gin-handler?path=* ; we restore the real path for Gin.
func Handler(w http.ResponseWriter, r *http.Request) {
	r2 := normalizeRequestForGin(r)
	ginHandler.ServeHTTP(w, r2)
}

func normalizeRequestForGin(r *http.Request) *http.Request {
	if p := r.URL.Query().Get("path"); p != "" {
		u := *r.URL
		u.Path = "/v1/" + strings.TrimPrefix(strings.TrimSpace(p), "/")
		q := r.URL.Query()
		q.Del("path")
		u.RawQuery = q.Encode()
		out := r.Clone(r.Context())
		out.URL = &u
		return out
	}

	// Fallback: path under /api/vercel/gin-handler/... -> /v1/...
	const prefix = "/api/vercel/gin-handler"
	if strings.HasPrefix(r.URL.Path, prefix) {
		rest := strings.TrimPrefix(r.URL.Path, prefix)
		rest = strings.TrimPrefix(rest, "/")
		if rest != "" {
			u := *r.URL
			u.Path = "/v1/" + rest
			out := r.Clone(r.Context())
			out.URL = &u
			return out
		}
	}

	return r
}
