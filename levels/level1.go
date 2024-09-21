package levels

import (
	"fmt"
	"strings"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
)

/*
	The goal of level1 is to create an nginx deployment
*/

type Level1 struct {
}

func (l *Level1) GetID() int {
	return 1
}

func (l *Level1) GetName() string {
	return "what is a pod"
}

func (l *Level1) GetDesiredCluster() client.Cluster {
	return client.Cluster{
		Deployments: []client.Deployment{
			{
				Name:     "nginx",
				Replicas: 1,
				Pods: []client.SimplePod{
					{
						Name:   "nginx",
						Mounts: []client.PodVolume{},
						Containers: []client.Container{
							{
								Name:  "nginx",
								Image: "nginx",
								Ports: []int32{80},
							},
						},
					},
				},
			},
		},
	}
}

func (l *Level1) GetClusterStatus(cluster client.Cluster, msg string) string {
	if len(cluster.Deployments) != 1 {
		return fmt.Sprintf("There should be 1 deployment; found: %d", len(cluster.Deployments))
	}

	deployment := cluster.Deployments[0]
	if len(deployment.Pods) != 1 {
		return fmt.Sprintf("There should be 1 pod; found: %d", len(deployment.Pods))
	}

	pod := deployment.Pods[0]
	if strings.HasPrefix(pod.Name, "nginx") == false {
		return "Pod was found but only named nginx; instead a pod named " + pod.Name + " was found"
	}

	if len(pod.Containers) != 1 {
		return fmt.Sprintf("There should be 1 container; found: %d", len(pod.Containers))
	}

	container := pod.Containers[0]

	if !(strings.Contains(container.Image, "nginx")) {
		return "Container image should be nginx; found: " + container.Image
	}

	if len(container.Ports) != 1 {
		return fmt.Sprintf("There should be 1 port; found: %d", len(container.Ports))
	}

	if container.Ports[0] != 80 {
		return fmt.Sprintf("Port should be 80; found: %d", container.Ports[0])
	}

	return "success"
}
