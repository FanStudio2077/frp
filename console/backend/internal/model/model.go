package model

import "time"

// Role values for RBAC.
const (
	RoleAdmin    = "admin"
	RoleOperator = "operator"
	RoleViewer   = "viewer"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	Disabled     bool      `json:"disabled"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Server struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Addr      string    `json:"addr"`
	Port      int       `json:"port"`
	TokenEnc  string    `json:"-"`
	TLSJSON   string    `json:"tls,omitempty"`
	Remark    string    `json:"remark,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type Node struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	TokenHash string    `json:"-"`
	IP        string    `json:"ip,omitempty"`
	OS        string    `json:"os,omitempty"`
	Arch      string    `json:"arch,omitempty"`
	Version   string    `json:"version,omitempty"`
	Online    bool      `json:"online"`
	LastSeen  time.Time `json:"lastSeen"`
	CreatedAt time.Time `json:"createdAt"`
}

type Instance struct {
	ID        int64  `json:"id"`
	NodeID    int64  `json:"nodeId"`
	ServerID  int64  `json:"serverId"`
	Name      string `json:"name"`
	AdminPort int    `json:"adminPort"`
	Status    string `json:"status"`
}

type Proxy struct {
	ID            int64  `json:"id"`
	InstanceID    int64  `json:"instanceId"`
	NodeID        int64  `json:"nodeId"`
	ServerID      int64  `json:"serverId"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	LocalIP       string `json:"localIp"`
	LocalPort     int    `json:"localPort"`
	RemotePort    int    `json:"remotePort,omitempty"`
	Subdomain     string `json:"subdomain,omitempty"`
	CustomDomains string `json:"customDomains,omitempty"`
	ExtraJSON     string `json:"extra,omitempty"`
	Enabled       bool   `json:"enabled"`
}

type AuditLog struct {
	ID        int64     `json:"id"`
	User      string    `json:"user"`
	Action    string    `json:"action"`
	Target    string    `json:"target"`
	Detail    string    `json:"detail,omitempty"`
	IP        string    `json:"ip,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type AlertRule struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Metric    string `json:"metric"`
	Threshold int64  `json:"threshold"`
	Scope     string `json:"scope"`
	Enabled   bool   `json:"enabled"`
}

type AlertEvent struct {
	ID        int64     `json:"id"`
	RuleID    int64     `json:"ruleId"`
	ServerID  int64     `json:"serverId"`
	Value     int64     `json:"value"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}
