package v1

import "testing"

func serversBoolPtr(b bool) *bool { return &b }

func TestClientCommonConfigCompleteServers(t *testing.T) {
	c := &ClientCommonConfig{}
	c.Servers = []ServerEndpoint{
		{ServerAddr: "1.1.1.1", ServerPort: 7000, Priority: 2},
		{ServerAddr: "2.2.2.2", Priority: 1},
		{ServerAddr: "3.3.3.3", ServerPort: 7001, Enabled: serversBoolPtr(false), Priority: 0},
	}
	if err := c.Complete(); err != nil {
		t.Fatal(err)
	}
	if c.Servers[0].Name != "server-1" || c.Servers[1].Name != "server-2" {
		t.Fatalf("names not defaulted: %#v", c.Servers)
	}
	if c.ServerAddr != "1.1.1.1" || c.ServerPort != 7000 {
		t.Fatalf("top-level not set from first enabled: %s:%d", c.ServerAddr, c.ServerPort)
	}
	if c.Servers[1].ServerPort != 7000 {
		t.Fatalf("server port default failed: %d", c.Servers[1].ServerPort)
	}
	if c.Servers[2].IsEnabled() {
		t.Fatal("server-3 should be disabled")
	}
}
