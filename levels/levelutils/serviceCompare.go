package levelutils

import (
	"fmt"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
)

func CompareServices(expected client.Service, current client.Service, expectedEndpoints []client.EndPoint) string {
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

	if len(expectedEndpoints) < 1 {
		return "success"
	}

	if len(expected.Endpoints) != len(current.Endpoints) {
		return fmt.Sprintf("Expected %d endpoints, but found %d for label '%s,%s' at service: %s", len(expected.Endpoints), len(current.Endpoints), current.SelectorMap["app"], current.SelectorMap["version"],
			current.Name)
	}
	return "success"
}
