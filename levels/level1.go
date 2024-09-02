package levels

import (
	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
	"strings"
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

func (l *Level1) GetClusterDiff(cluster client.Cluster) string {
	if len(cluster.Deployments) != 1 {
		return "There should be 1 deployment; found 	" + string(rune(len(cluster.Deployments)))
	}

	deployment := cluster.Deployments[0]
	if len(deployment.Pods) != 1 {
		return "There should be 1 pod; found:" + string(rune(len(deployment.Pods)))
	}

	pod := deployment.Pods[0]
	if pod.Name != "nginx" {
		return "Pod was found but only named nginx; instead a pod named " + pod.Name + " was found"
	}

	if len(pod.Containers) != 1 {
		return "There should be 1 container; found: " + string(rune(len(pod.Containers)))
	}

	container := pod.Containers[0]
	if container.Name != "nginx" {
		return "Container was found but only named nginx; instead a container named " + container.Name + " was found"
	}

	if !(strings.Contains(container.Image, "nginx:")) {
		return "Container image should be nginx; found: " + container.Image
	}

	if len(container.Ports) != 1 {
		return "There should be 1 port; found: " + string(rune(len(container.Ports)))
	}

	if container.Ports[0] != 80 {
		return "Container port should be 80; found: " + string(rune(container.Ports[0]))
	}

	return ""
}
