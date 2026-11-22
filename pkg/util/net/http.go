// Copyright 2017 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package net

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fatedier/frp/pkg/util/util"
)

type HTTPAuthMiddleware struct {
	user          string
	passwd        string
	authFailDelay time.Duration
	users         map[string]string
	roles         map[string]string
	mu            sync.RWMutex
}

type authRoleKey struct{}

// RoleFromContext returns the authenticated user role stored by the auth middleware.
func RoleFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(authRoleKey{}).(string); ok {
		return v
	}
	return ""
}

func NewHTTPAuthMiddleware(user, passwd string) *HTTPAuthMiddleware {
	return &HTTPAuthMiddleware{
		user:   user,
		passwd: passwd,
	}
}

func (authMid *HTTPAuthMiddleware) SetAuthFailDelay(delay time.Duration) *HTTPAuthMiddleware {
	authMid.authFailDelay = delay
	return authMid
}

// SetUsers configures additional users and their roles.
func (authMid *HTTPAuthMiddleware) SetUsers(users map[string]string, roles map[string]string) *HTTPAuthMiddleware {
	authMid.mu.Lock()
	authMid.users = users
	authMid.roles = roles
	authMid.mu.Unlock()
	return authMid
}

func (authMid *HTTPAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authMid.mu.RLock()
		legacyUser, legacyPass := authMid.user, authMid.passwd
		users, roles := authMid.users, authMid.roles
		authMid.mu.RUnlock()
		reqUser, reqPasswd, hasAuth := r.BasicAuth()
		if legacyUser == "" && legacyPass == "" && len(users) == 0 {
			next.ServeHTTP(w, r)
			return
		}
		if hasAuth {
			role := ""
			if legacyUser != "" || legacyPass != "" {
				if util.ConstantTimeEqString(reqUser, legacyUser) && util.ConstantTimeEqString(reqPasswd, legacyPass) {
					role = "admin"
				}
			}
			if role == "" {
				if pw, ok := users[reqUser]; ok && util.ConstantTimeEqString(reqPasswd, pw) {
					role = roles[reqUser]
					if role == "" {
						role = "viewer"
					}
				}
			}
			if role != "" {
				ctx := context.WithValue(r.Context(), authRoleKey{}, role)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		if authMid.authFailDelay > 0 {
			time.Sleep(authMid.authFailDelay)
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	})
}

type HTTPGzipWrapper struct {
	h http.Handler
}

func (gw *HTTPGzipWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		gw.h.ServeHTTP(w, r)
		return
	}
	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	defer gz.Close()
	gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
	gw.h.ServeHTTP(gzr, r)
}

func MakeHTTPGzipHandler(h http.Handler) http.Handler {
	return &HTTPGzipWrapper{
		h: h,
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
