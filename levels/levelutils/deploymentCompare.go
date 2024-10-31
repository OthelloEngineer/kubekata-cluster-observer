package levelutils

import (
	"fmt"
	"strings"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
)

func CompareDeployments(currentDeployment client.Deployment, expectedDeployment client.Deployment) string {

	if currentDeployment.Name != expectedDeployment.Name {
		return fmt.Sprintf("Deployment name should be %s; found: %s", expectedDeployment.Name, currentDeployment.Name)
	}

	if currentDeployment.Replicas != expectedDeployment.Replicas {
		return fmt.Sprintf("Deployment %s should have %d replicas; found: %d", currentDeployment.Name, expectedDeployment.Replicas, currentDeployment.Replicas)
	}

	if len(currentDeployment.Pods) != len(expectedDeployment.Pods) {
		return fmt.Sprintf("Deployment %s should have %d pods; found: %d", currentDeployment.Name, len(expectedDeployment.Pods), len(currentDeployment.Pods))
	}

	for j, expectedPod := range expectedDeployment.Pods {
		currentPod := currentDeployment.Pods[j]
		// Expected pod name should be in the current pod name
		if !strings.Contains(currentPod.Name, expectedPod.Name) {
			return fmt.Sprintf("Pod name should be %s; found: %s", expectedPod.Name, currentPod.Name)
		}

		if len(currentPod.Containers) != len(expectedPod.Containers) {
			return fmt.Sprintf("Pod %s should have %d containers; found: %d", currentPod.Name, len(expectedPod.Containers), len(currentPod.Containers))
		}

		for k, expectedContainer := range expectedPod.Containers {
			currentContainer := currentPod.Containers[k]
			if currentContainer.Name != expectedContainer.Name {
				return fmt.Sprintf("Container name should be %s; found: %s", expectedContainer.Name, currentContainer.Name)
			}

			if currentContainer.Image != expectedContainer.Image {
				return fmt.Sprintf("Container %s should have image %s; found: %s", currentContainer.Name, expectedContainer.Image, currentContainer.Image)
			}

			if len(currentContainer.Ports) != len(expectedContainer.Ports) {
				return fmt.Sprintf("Container %s should have %d ports; found: %d", currentContainer.Name, len(expectedContainer.Ports), len(currentContainer.Ports))
			}

			for l, expectedPort := range expectedContainer.Ports {
				currentPort := currentContainer.Ports[l]
				if currentPort != expectedPort {
					return fmt.Sprintf("Container %s should have port %d; found: %d", currentContainer.Name, expectedPort, currentPort)
				}
			}
			if len(currentContainer.Mount) != 0 && len(expectedContainer.Mount) == 0 {
				if len(currentContainer.Mount) != len(expectedContainer.Mount) {
					return fmt.Sprintf("Container %s should have %d mounts; found: %d", currentContainer.Name, len(expectedContainer.Mount), len(currentContainer.Mount))
				}

				for m, expectedMount := range expectedContainer.Mount {
					currentMount := currentContainer.Mount[m]
					if currentMount != expectedMount {
						return fmt.Sprintf("Mount name should be %s; found: %s", expectedMount, currentMount)
					}
				}
			}
		}
	}
	return "success"
}

func ComparePods(currentPods []client.SimplePod, expectedPods []client.SimplePod) string {
	if len(currentPods) != len(expectedPods) {
		return fmt.Sprintf("There should be %d pods; found: %d", len(expectedPods), len(currentPods))
	}

	for i, expectedPod := range expectedPods {
		currentPod := currentPods[i]

		if len(currentPod.Containers) != len(expectedPod.Containers) {
			return fmt.Sprintf("Pod %s should have %d containers; found: %d", currentPod.Name, len(expectedPod.Containers), len(currentPod.Containers))
		}

		CompareContainers(currentPod.Containers, expectedPod.Containers)
	}
	return "success"
}

// CompareContainers compares the current containers with the expected containers in terms of
// name, image, ports, mounts, requests, and limits. In cases where limits are not set, it should be notes at ""
func CompareContainers(currentContainers []client.Container, expectedContainers []client.Container) string {
	if len(currentContainers) != len(expectedContainers) {
		return fmt.Sprintf("There should be %d containers; found: %d", len(expectedContainers), len(currentContainers))
	}

	for i, expectedContainer := range expectedContainers {
		currentContainer := currentContainers[i]
		if currentContainer.Name != expectedContainer.Name {
			return fmt.Sprintf("Container name should be %s; found: %s", expectedContainer.Name, currentContainer.Name)
		}

		if currentContainer.Image != expectedContainer.Image {
			return fmt.Sprintf("Container %s should have image %s; found: %s", currentContainer.Name, expectedContainer.Image, currentContainer.Image)
		}

		if len(currentContainer.Ports) != len(expectedContainer.Ports) {
			return fmt.Sprintf("Container %s should have %d ports; found: %d", currentContainer.Name, len(expectedContainer.Ports), len(currentContainer.Ports))
		}

		if len(currentContainer.Mount) != 0 && len(expectedContainer.Mount) == 0 {
			if len(currentContainer.Mount) != len(expectedContainer.Mount) {
				return fmt.Sprintf("Container %s should have %d mounts; found: %d", currentContainer.Name, len(expectedContainer.Mount), len(currentContainer.Mount))
			}
		}

		msg := CompareResources(currentContainer.Requests, expectedContainer.Requests)
		if msg != "success" {
			return msg
		}

		msg = CompareResources(currentContainer.Limits, expectedContainer.Limits)
		if msg != "success" {
			return msg
		}

		for j, expectedPort := range expectedContainer.Ports {
			currentPort := currentContainer.Ports[j]
			if currentPort != expectedPort {
				return fmt.Sprintf("Container %s should have port %d; found: %d", currentContainer.Name, expectedPort, currentPort)
			}
		}
	}
	return "success"
}

func CompareResources(currentResource client.Resource, expectedResource client.Resource) string {
	if expectedResource.CPU == "" && expectedResource.Memory == "" {
		return "success"
	}
	if currentResource.CPU != expectedResource.CPU {
		return fmt.Sprintf("CPU should be %s; found: %s", expectedResource.CPU, currentResource.CPU)
	}

	if currentResource.Memory != expectedResource.Memory {
		return fmt.Sprintf("Memory should be %s; found: %s", expectedResource.Memory, currentResource.Memory)
	}

	return "success"
}

func CompareImagesAndPortOfDeployments(currentDeployments []client.Deployment, expectedDeployments []client.Deployment) string {
	if len(currentDeployments) != len(expectedDeployments) {
		return fmt.Sprintf("There should be %d deployments; found: %d", len(expectedDeployments), len(currentDeployments))
	}

	expectedImagePort := map[string][]int32{}
	currentImagePort := map[string][]int32{}

	for _, expectedDeployment := range expectedDeployments {
		pod := expectedDeployment.Pods[0]
		for _, container := range pod.Containers {
			expectedImagePort[container.Image] = container.Ports
			fmt.Println("Expected Image Port: ", expectedImagePort[container.Image])
		}
	}

	for _, currentDeployment := range currentDeployments {
		if len(currentDeployment.Pods) == 0 {
			return fmt.Sprintf("Deployment %s should have at least 1 pod; found: 0", currentDeployment.Name)
		}
		pod := currentDeployment.Pods[0]
		for _, container := range pod.Containers {
			currentImagePort[container.Image] = container.Ports
		}
	}

	for image, ports := range expectedImagePort {
		currentPorts, ok := currentImagePort[image]
		if !ok {
			return fmt.Sprintf("Container image %s not found", image)
		}
		for i, port := range ports {
			if currentPorts[i] != port {
				return fmt.Sprintf("Container image %s should have port %d; found: %d", image, port, currentPorts[i])
			}
		}
	}
	return "success"
}
