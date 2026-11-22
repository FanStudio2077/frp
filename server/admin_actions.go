package server

import "fmt"

// DisconnectClient forcefully closes the control connection of the client with the given registry key.
func (svr *Service) DisconnectClient(key string) error {
	info, ok := svr.clientRegistry.GetByKey(key)
	if !ok {
		return fmt.Errorf("client %q not found", key)
	}
	ctl, ok := svr.ctlManager.GetByID(info.RunID)
	if !ok || ctl == nil {
		return fmt.Errorf("client %q is not online", key)
	}
	return ctl.Close()
}

// DeleteProxiesByName closes and removes the named proxies. Returns the number removed.
func (svr *Service) DeleteProxiesByName(names []string) int {
	removed := 0
	for _, name := range names {
		if pxy, ok := svr.pxyManager.GetByName(name); ok {
			pxy.Close()
			svr.pxyManager.Del(name)
			removed++
		}
	}
	return removed
}
