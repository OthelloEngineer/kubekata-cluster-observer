package levelutils

import (
	"fmt"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
)

func compareServices(expected client.Service, current client.Service) string {
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
		if found == false {
			return fmt.Sprintf("Expected service port could not find: %s ", port)
		}
	}

	for key, value := range expected.Selector {
		if current.Selector[key] != value {
			return fmt.Sprintf("Expected service selector could not find: %s ", key)
		}
	}
	return "success"
}
