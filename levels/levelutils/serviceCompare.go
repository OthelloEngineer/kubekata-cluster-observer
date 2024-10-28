package levelutils

import (
	"fmt"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
)

func CompareServices(expected client.Service, current client.Service) string {
	if expected.Name != current.Name {
		return "Service name is not correct"
	}
	if expected.Type != current.Type {
		return "Service type is not correct"
	}

	for _, port := range expected.Ports {
		found := false
		for _, currentPort := range current.Ports {
			if port == currentPort {
				found = true
				break
			}
		}
		if !found {
			return fmt.Sprintf("Expected service port could not find: %d ", port)
		}
	}

	for key, value := range expected.SelectorMap {
		if current.SelectorMap[key] != value {
			return fmt.Sprintf("Expected service selector could not find: %s ", key)
		}
	}

	if len(expected.Endpoints) < 1 {
		return "success"
	}

	expectedEndpointCount := len(expected.Endpoints)
	currentEndpointCount := len(current.Endpoints)
	if expectedEndpointCount != currentEndpointCount {
		return fmt.Sprintf("Expected %d points, but the service is connected to %d pods ", expectedEndpointCount, currentEndpointCount)
	}

	return "success"
}
