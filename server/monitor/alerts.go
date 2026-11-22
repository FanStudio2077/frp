package monitor

import (
	"sync"
	"time"

	v1 "github.com/fatedier/frp/pkg/config/v1"
)

// Alert is one fired monitoring alert.
type Alert struct {
	Time      time.Time `json:"time"`
	Level     string    `json:"level"`
	Metric    string    `json:"metric"`
	Value     int64     `json:"value"`
	Threshold int64     `json:"threshold"`
	Message   string    `json:"message"`
}

// Engine evaluates monitoring thresholds and keeps a bounded alert history.
type Engine struct {
	mu       sync.Mutex
	cfg      v1.AlertConfig
	alerts   []Alert
	max      int
	lastFire map[string]time.Time
}

// NewEngine creates an alert engine.
func NewEngine(cfg v1.AlertConfig) *Engine {
	return &Engine{cfg: cfg, max: 100, lastFire: map[string]time.Time{}}
}

// UpdateConfig replaces the thresholds at runtime.
func (e *Engine) UpdateConfig(cfg v1.AlertConfig) {
	e.mu.Lock()
	e.cfg = cfg
	e.mu.Unlock()
}

// Config returns the current alert configuration.
func (e *Engine) Config() v1.AlertConfig {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.cfg
}

// Evaluate checks thresholds and records alerts (60s cooldown per metric).
func (e *Engine) Evaluate(now time.Time, onlineClients, curConns int64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.cfg.Enabled {
		return
	}
	e.checkLocked(now, "onlineClients", onlineClients, e.cfg.MaxOnlineClients)
	e.checkLocked(now, "curConns", curConns, e.cfg.MaxCurConns)
}

func (e *Engine) checkLocked(now time.Time, metric string, value, threshold int64) {
	if threshold <= 0 || value <= threshold {
		return
	}
	if t, ok := e.lastFire[metric]; ok && now.Sub(t) < time.Minute {
		return
	}
	e.lastFire[metric] = now
	e.alerts = append(e.alerts, Alert{
		Time: now, Level: "warning", Metric: metric, Value: value, Threshold: threshold,
		Message: metric + " exceeded threshold",
	})
	if len(e.alerts) > e.max {
		e.alerts = e.alerts[len(e.alerts)-e.max:]
	}
}

// List returns the recorded alerts (oldest first).
func (e *Engine) List() []Alert {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]Alert, len(e.alerts))
	copy(out, e.alerts)
	return out
}

// Clear removes all recorded alerts.
func (e *Engine) Clear() {
	e.mu.Lock()
	e.alerts = nil
	e.lastFire = map[string]time.Time{}
	e.mu.Unlock()
}
