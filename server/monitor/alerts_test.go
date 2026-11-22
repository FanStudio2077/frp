package monitor

import (
	"testing"
	"time"

	v1 "github.com/fatedier/frp/pkg/config/v1"
)

func TestEngineEvaluate(t *testing.T) {
	e := NewEngine(v1.AlertConfig{Enabled: true, MaxOnlineClients: 2, MaxCurConns: 10})
	now := time.Now()
	e.Evaluate(now, 1, 5)
	if len(e.List()) != 0 {
		t.Fatal("no alert expected below thresholds")
	}
	e.Evaluate(now, 3, 5)
	items := e.List()
	if len(items) != 1 || items[0].Metric != "onlineClients" {
		t.Fatalf("expected onlineClients alert, got %+v", items)
	}
	e.Evaluate(now.Add(30*time.Second), 4, 5)
	if len(e.List()) != 1 {
		t.Fatal("cooldown should suppress duplicate alert")
	}
	e.Evaluate(now.Add(2*time.Minute), 4, 11)
	if len(e.List()) != 3 {
		t.Fatalf("expected 3 alerts, got %d", len(e.List()))
	}
	e.Clear()
	if len(e.List()) != 0 {
		t.Fatal("clear failed")
	}
}
