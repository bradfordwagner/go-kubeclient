package kube

import (
	"fmt"
)

// GetServerInfo returns the Kubernetes API server URL and CA certificate data
// Returns an error if the rest.Config is not available
func (c *client) GetServerInfo() (server string, ca string, err error) {
	if c.config == nil {
		err = fmt.Errorf("rest config is not available - use NewClientInterfaceWithConfig or NewDefaultClientInterface")
		return
	}

	server = c.config.Host
	if c.config.TLSClientConfig.CAData != nil {
		ca = string(c.config.TLSClientConfig.CAData)
	}

	return
}
