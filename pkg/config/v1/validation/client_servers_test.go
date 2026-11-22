package validation

import (
	"testing"

	v1 "github.com/fatedier/frp/pkg/config/v1"
)

func TestValidateServers(t *testing.T) {
	dup := &v1.ClientCommonConfig{Servers: []v1.ServerEndpoint{
		{Name: "a", ServerAddr: "1.1.1.1", ServerPort: 7000},
		{Name: "a", ServerAddr: "2.2.2.2", ServerPort: 7000},
	}}
	if _, err := validateServers(dup); err == nil {
		t.Fatal("expected duplicate name error")
	}
	bad := &v1.ClientCommonConfig{Servers: []v1.ServerEndpoint{
		{Name: "a", ServerAddr: "", ServerPort: 0},
	}}
	if _, err := validateServers(bad); err == nil {
		t.Fatal("expected addr/port error")
	}
	allDisabled := false
	_ = allDisabled
	disabled := &v1.ClientCommonConfig{Servers: []v1.ServerEndpoint{
		{Name: "a", ServerAddr: "1.1.1.1", ServerPort: 7000, Enabled: serversBoolPtr2(false)},
	}}
	if _, err := validateServers(disabled); err == nil {
		t.Fatal("expected at least one enabled error")
	}
	ok := &v1.ClientCommonConfig{Servers: []v1.ServerEndpoint{
		{Name: "a", ServerAddr: "1.1.1.1", ServerPort: 7000},
	}}
	if _, err := validateServers(ok); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func serversBoolPtr2(b bool) *bool { return &b }
