package v1

// AlertConfig configures runtime monitoring thresholds for frps (hot-reloadable).
type AlertConfig struct {
	Enabled              bool  `json:"enabled,omitempty"`
	MaxOnlineClients     int64 `json:"maxOnlineClients,omitempty"`
	MaxCurConns          int64 `json:"maxCurConns,omitempty"`
	CheckIntervalSeconds int   `json:"checkIntervalSeconds,omitempty"`
}
