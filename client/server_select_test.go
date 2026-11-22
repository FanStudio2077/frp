package client

import (
	"testing"

	v1 "github.com/fatedier/frp/pkg/config/v1"
)

func TestBuildServerEntriesOrdering(t *testing.T) {
	list := []v1.ServerEndpoint{
		{Name: "b", ServerAddr: "2.2.2.2", ServerPort: 7000, Priority: 2},
		{Name: "a", ServerAddr: "1.1.1.1", ServerPort: 7000, Priority: 1},
	}
	es := buildServerEntries(list, "", 0)
	if len(es) != 2 || es[0].Name != "a" || es[1].Name != "b" {
		t.Fatalf("unexpected order: %#v", es)
	}
}

func TestPickServerCycles(t *testing.T) {
	svr := &Service{common: &v1.ClientCommonConfig{}}
	svr.servers = []serverEntry{
		{Name: "a", ServerAddr: "1.1.1.1", ServerPort: 7000, Enabled: true},
		{Name: "b", ServerAddr: "2.2.2.2", ServerPort: 7001, Enabled: true},
	}
	svr.pickServer()
	if svr.ActiveServer() != "a" || svr.common.ServerAddr != "1.1.1.1" {
		t.Fatalf("first active=%s addr=%s", svr.ActiveServer(), svr.common.ServerAddr)
	}
	svr.pickServer()
	if svr.ActiveServer() != "b" || svr.common.ServerPort != 7001 {
		t.Fatalf("second active=%s port=%d", svr.ActiveServer(), svr.common.ServerPort)
	}
	svr.pickServer()
	if svr.ActiveServer() != "a" {
		t.Fatalf("third active=%s", svr.ActiveServer())
	}
}
