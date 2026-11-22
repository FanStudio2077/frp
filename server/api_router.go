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

package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	httppkg "github.com/fatedier/frp/pkg/util/http"
	netpkg "github.com/fatedier/frp/pkg/util/net"
	"github.com/fatedier/frp/server/audit"
	adminapi "github.com/fatedier/frp/server/http"
)

func (svr *Service) registerRouteHandlers(helper *httppkg.RouterRegisterHelper) {
	helper.Router.HandleFunc("/healthz", healthz)
	subRouter := helper.Router.NewRoute().Subrouter()

	subRouter.Use(helper.AuthMiddleware)
	subRouter.Use(requireAdminForWrites)
	subRouter.Use(httppkg.NewRequestLogger)
	subRouter.Use(auditMiddleware(svr.auditStore))

	// metrics
	if svr.cfg.EnablePrometheus {
		subRouter.Handle("/metrics", promhttp.Handler())
	}

	apiController := adminapi.NewController(svr.cfg, svr.clientRegistry, svr.pxyManager)

	// apis
	subRouter.HandleFunc("/api/serverinfo", httppkg.MakeHTTPHandlerFunc(apiController.APIServerInfo)).Methods("GET")
	subRouter.HandleFunc("/api/proxy/{type}", httppkg.MakeHTTPHandlerFunc(apiController.APIProxyByType)).Methods("GET")
	subRouter.HandleFunc("/api/proxy/{type}/{name}", httppkg.MakeHTTPHandlerFunc(apiController.APIProxyByTypeAndName)).Methods("GET")
	subRouter.HandleFunc("/api/proxies/{name}", httppkg.MakeHTTPHandlerFunc(apiController.APIProxyByName)).Methods("GET")
	subRouter.HandleFunc("/api/traffic/{name}", httppkg.MakeHTTPHandlerFunc(apiController.APIProxyTraffic)).Methods("GET")
	subRouter.HandleFunc("/api/clients", httppkg.MakeHTTPHandlerFunc(apiController.APIClientList)).Methods("GET")
	subRouter.HandleFunc("/api/clients/{key}", httppkg.MakeHTTPHandlerFunc(apiController.APIClientDetail)).Methods("GET")
	subRouter.HandleFunc("/api/proxies", httppkg.MakeHTTPHandlerFunc(apiController.DeleteProxies)).Methods("DELETE")

	subRouter.HandleFunc("/api/v2/users", httppkg.MakeHTTPHandlerFuncV2(apiController.APIV2UserList)).Methods("GET")
	subRouter.HandleFunc("/api/v2/system/info", httppkg.MakeHTTPHandlerFuncV2(apiController.APIV2SystemInfo)).Methods("GET")
	subRouter.HandleFunc("/api/v2/system/prune", httppkg.MakeHTTPHandlerFuncV2(apiController.APIV2SystemPrune)).Methods("POST")
	subRouter.HandleFunc("/api/v2/clients", httppkg.MakeHTTPHandlerFuncV2(apiController.APIV2ClientList)).Methods("GET")
	v2EncodedPathRouter := subRouter.NewRoute().Subrouter()
	v2EncodedPathRouter.UseEncodedPath()
	v2EncodedPathRouter.HandleFunc("/api/v2/clients/{key}", httppkg.MakeHTTPHandlerFuncV2(apiController.APIV2ClientDetail)).Methods("GET")
	subRouter.HandleFunc("/api/v2/proxies", httppkg.MakeHTTPHandlerFuncV2(apiController.APIV2ProxyList)).Methods("GET")
	v2EncodedPathRouter.HandleFunc("/api/v2/proxies/{name}/traffic", httppkg.MakeHTTPHandlerFuncV2(apiController.APIV2ProxyTraffic)).Methods("GET")
	v2EncodedPathRouter.HandleFunc("/api/v2/proxies/{name}", httppkg.MakeHTTPHandlerFuncV2(apiController.APIV2ProxyDetail)).Methods("GET")

	v2EncodedPathRouter.HandleFunc("/api/v2/clients/{key}/disconnect", httppkg.MakeHTTPHandlerFuncV2(func(ctx *httppkg.Context) (any, error) {
		if err := svr.DisconnectClient(ctx.Param("key")); err != nil {
			return nil, httppkg.NewError(http.StatusNotFound, err.Error())
		}
		return nil, nil
	})).Methods("POST")
	subRouter.HandleFunc("/api/v2/proxies", httppkg.MakeHTTPHandlerFuncV2(func(ctx *httppkg.Context) (any, error) {
		body, err := ctx.Body()
		if err != nil {
			return nil, httppkg.NewError(http.StatusBadRequest, err.Error())
		}
		var payload struct {
			Names []string `json:"names"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, httppkg.NewError(http.StatusBadRequest, err.Error())
		}
		return map[string]any{"deleted": svr.DeleteProxiesByName(payload.Names)}, nil
	})).Methods("DELETE")

	// config reload + web dashboard users (RBAC)
	subRouter.HandleFunc("/api/v2/system/reload", httppkg.MakeHTTPHandlerFuncV2(func(ctx *httppkg.Context) (any, error) {
		res, err := svr.ReloadConfigFromFile()
		if err != nil {
			return nil, httppkg.NewError(http.StatusBadRequest, err.Error())
		}
		return res, nil
	})).Methods("POST")
	subRouter.HandleFunc("/api/v2/web/users", httppkg.MakeHTTPHandlerFuncV2(func(ctx *httppkg.Context) (any, error) {
		type webUserView struct {
			User   string `json:"user"`
			Role   string `json:"role"`
			Legacy bool   `json:"legacy,omitempty"`
		}
		items := make([]webUserView, 0)
		if svr.cfg.WebServer.User != "" {
			items = append(items, webUserView{User: svr.cfg.WebServer.User, Role: "admin", Legacy: true})
		}
		for _, u := range svr.cfg.WebServer.Users {
			items = append(items, webUserView{User: u.User, Role: u.Role})
		}
		return map[string]any{"items": items}, nil
	})).Methods("GET")
	subRouter.HandleFunc("/api/v2/web/users", httppkg.MakeHTTPHandlerFuncV2(func(ctx *httppkg.Context) (any, error) {
		body, err := ctx.Body()
		if err != nil {
			return nil, httppkg.NewError(http.StatusBadRequest, err.Error())
		}
		var payload struct {
			Users []v1.WebUser `json:"users"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, httppkg.NewError(http.StatusBadRequest, err.Error())
		}
		if err := svr.UpdateWebUsers(payload.Users); err != nil {
			return nil, httppkg.NewError(http.StatusBadRequest, err.Error())
		}
		return nil, nil
	})).Methods("PUT")

	// monitoring alerts
	subRouter.HandleFunc("/api/v2/alerts", httppkg.MakeHTTPHandlerFuncV2(func(ctx *httppkg.Context) (any, error) {
		cfg := svr.alertEngine.Config()
		return map[string]any{"enabled": cfg.Enabled, "thresholds": cfg, "items": svr.alertEngine.List()}, nil
	})).Methods("GET")
	subRouter.HandleFunc("/api/v2/alerts", httppkg.MakeHTTPHandlerFuncV2(func(ctx *httppkg.Context) (any, error) {
		svr.alertEngine.Clear()
		return nil, nil
	})).Methods("DELETE")
	subRouter.HandleFunc("/api/v2/system/alerts", httppkg.MakeHTTPHandlerFuncV2(func(ctx *httppkg.Context) (any, error) {
		body, err := ctx.Body()
		if err != nil {
			return nil, httppkg.NewError(http.StatusBadRequest, err.Error())
		}
		var cfg v1.AlertConfig
		if err := json.Unmarshal(body, &cfg); err != nil {
			return nil, httppkg.NewError(http.StatusBadRequest, err.Error())
		}
		svr.alertEngine.UpdateConfig(cfg)
		svr.cfg.Alerts = cfg
		return nil, nil
	})).Methods("PUT")
	subRouter.HandleFunc("/api/v2/system/config", httppkg.MakeHTTPHandlerFuncV2(func(ctx *httppkg.Context) (any, error) {
		c := svr.cfg
		return map[string]any{
			"bindAddr": c.BindAddr, "bindPort": c.BindPort, "kcpBindPort": c.KCPBindPort,
			"quicBindPort": c.QUICBindPort, "vhostHTTPPort": c.VhostHTTPPort, "vhostHTTPSPort": c.VhostHTTPSPort,
			"tcpmuxHTTPConnectPort": c.TCPMuxHTTPConnectPort, "subdomainHost": c.SubDomainHost,
			"maxPortsPerClient": c.MaxPortsPerClient, "alerts": c.Alerts,
		}, nil
	})).Methods("GET")

	// audit log
	subRouter.HandleFunc("/api/v2/audit/logs", httppkg.MakeHTTPHandlerFuncV2(func(ctx *httppkg.Context) (any, error) {
		limit, _ := strconv.Atoi(ctx.Query("limit"))
		offset, _ := strconv.Atoi(ctx.Query("offset"))
		items, total := svr.auditStore.List(limit, offset, ctx.Query("method"))
		if items == nil {
			items = []audit.Entry{}
		}
		return map[string]any{"total": total, "items": items}, nil
	})).Methods("GET")
	subRouter.HandleFunc("/api/v2/audit/logs", httppkg.MakeHTTPHandlerFuncV2(func(ctx *httppkg.Context) (any, error) {
		svr.auditStore.Clear()
		return nil, nil
	})).Methods("DELETE")

	// view
	subRouter.Handle("/favicon.ico", http.FileServer(helper.AssetsFS)).Methods("GET")
	subRouter.PathPrefix("/static/").Handler(
		netpkg.MakeHTTPGzipHandler(http.StripPrefix("/static/", http.FileServer(helper.AssetsFS))),
	).Methods("GET")

	subRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/", http.StatusMovedPermanently)
	})
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
}

type auditStatusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *auditStatusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// auditMiddleware records mutating admin operations into the audit store.
func requireAdminForWrites(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}
		if netpkg.RoleFromContext(r.Context()) == "admin" {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "forbidden: admin role required", http.StatusForbidden)
	})
}

func auditMiddleware(store *audit.Store) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
			default:
				next.ServeHTTP(w, r)
				return
			}
			rec := &auditStatusRecorder{ResponseWriter: w, status: 200}
			next.ServeHTTP(rec, r)
			user := ""
			if u, _, ok := r.BasicAuth(); ok {
				user = u
			}
			store.Add(audit.Entry{Time: time.Now(), User: user, Method: r.Method, Path: r.URL.Path, Status: rec.status})
		})
	}
}
