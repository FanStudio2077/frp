package server

import (
	"fmt"
	"os"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/fatedier/frp/pkg/auth"
	"github.com/fatedier/frp/pkg/config"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/util/log"
)

// ReloadResult reports which settings were applied live and which need a restart.
type ReloadResult struct {
	Applied         []string `json:"applied"`
	RestartRequired []string `json:"restartRequired"`
}

// SetConfigFilePath records the config file used to start frps.
func (svr *Service) SetConfigFilePath(path string) { svr.configFilePath = path }

// ReloadConfigFromFile re-reads the config file and applies hot-reloadable settings.
func (svr *Service) ReloadConfigFromFile() (*ReloadResult, error) {
	if svr.configFilePath == "" {
		return nil, fmt.Errorf("frps was started without a config file; hot reload is unavailable")
	}
	newCfg, _, err := config.LoadServerConfig(svr.configFilePath, false)
	if err != nil {
		return nil, err
	}
	old := svr.cfg
	res := &ReloadResult{}
	addIf := func(name string, changed bool) {
		if changed {
			res.RestartRequired = append(res.RestartRequired, name)
		}
	}
	addIf("bindAddr", old.BindAddr != newCfg.BindAddr)
	addIf("bindPort", old.BindPort != newCfg.BindPort)
	addIf("kcpBindPort", old.KCPBindPort != newCfg.KCPBindPort)
	addIf("quicBindPort", old.QUICBindPort != newCfg.QUICBindPort)
	addIf("proxyBindAddr", old.ProxyBindAddr != newCfg.ProxyBindAddr)
	addIf("vhostHTTPPort", old.VhostHTTPPort != newCfg.VhostHTTPPort)
	addIf("vhostHTTPSPort", old.VhostHTTPSPort != newCfg.VhostHTTPSPort)
	addIf("tcpmuxHTTPConnectPort", old.TCPMuxHTTPConnectPort != newCfg.TCPMuxHTTPConnectPort)
	addIf("subDomainHost", old.SubDomainHost != newCfg.SubDomainHost)
	addIf("webServer.addr/port/tls", old.WebServer.Addr != newCfg.WebServer.Addr || old.WebServer.Port != newCfg.WebServer.Port || (old.WebServer.TLS == nil) != (newCfg.WebServer.TLS == nil))
	addIf("metricsServer", old.MetricsServer.Addr != newCfg.MetricsServer.Addr || old.MetricsServer.Port != newCfg.MetricsServer.Port)

	newAuth, err := auth.BuildServerAuth(&newCfg.Auth)
	if err != nil {
		return nil, err
	}

	svr.cfgMu.Lock()
	svr.cfg.Auth = newCfg.Auth
	svr.cfg.Alerts = newCfg.Alerts
	svr.cfg.DetailedErrorsToClient = newCfg.DetailedErrorsToClient
	svr.cfg.MaxPortsPerClient = newCfg.MaxPortsPerClient
	svr.cfg.UserConnTimeout = newCfg.UserConnTimeout
	svr.cfg.WebServer.Users = newCfg.WebServer.Users
	svr.cfgMu.Unlock()
	svr.auth = newAuth
	if svr.alertEngine != nil {
		svr.alertEngine.UpdateConfig(newCfg.Alerts)
	}
	if svr.webServer != nil {
		svr.webServer.SetUsers(newCfg.WebServer.Users)
	}
	res.Applied = []string{"auth", "alerts", "detailedErrorsToClient", "maxPortsPerClient", "userConnTimeout", "webServer.users"}
	log.Infof("frps config reloaded from %s; applied=%v restartRequired=%v", svr.configFilePath, res.Applied, res.RestartRequired)
	return res, nil
}

// UpdateWebUsers replaces dashboard users, persists them to the config file, and applies them live.
func (svr *Service) UpdateWebUsers(users []v1.WebUser) error {
	existing := map[string]v1.WebUser{}
	svr.cfgMu.RLock()
	for _, u := range svr.cfg.WebServer.Users {
		existing[u.User] = u
	}
	svr.cfgMu.RUnlock()
	for i := range users {
		if users[i].Password == "" {
			if old, ok := existing[users[i].User]; ok {
				users[i].Password = old.Password
			}
		}
		if users[i].Role == "" {
			users[i].Role = "viewer"
		}
	}
	if svr.configFilePath != "" {
		if err := writeWebUsersToFile(svr.configFilePath, users); err != nil {
			return err
		}
	}
	svr.cfgMu.Lock()
	svr.cfg.WebServer.Users = users
	svr.cfgMu.Unlock()
	if svr.webServer != nil {
		svr.webServer.SetUsers(users)
	}
	return nil
}

func writeWebUsersToFile(path string, users []v1.WebUser) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var root map[string]any
	if err := toml.Unmarshal(data, &root); err != nil {
		return err
	}
	ws, ok := root["webServer"].(map[string]any)
	if !ok {
		ws = map[string]any{}
		root["webServer"] = ws
	}
	arr := make([]map[string]any, 0, len(users))
	for _, u := range users {
		arr = append(arr, map[string]any{"user": u.User, "password": u.Password, "role": u.Role})
	}
	ws["users"] = arr
	out, err := toml.Marshal(root)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}
