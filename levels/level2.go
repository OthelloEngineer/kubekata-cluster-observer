package levels

import (
	"strings"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
	"k8s.io/apimachinery/pkg/labels"
)

/*
	The goal of level1 is to create an nginx deployment
*/

type Level2 struct {
}

func (l *Level2) GetID() int {
	return 2
}

func (l *Level2) GetName() string {
	return "what is a pod"
}

func (l *Level2) GetDesiredCluster() client.Cluster {
	return client.Cluster{
		Deployments: []client.Deployment{
			{
				Name:     "nginx-pod",
				Replicas: 1,
				SelectorMap: labels.Set{
					"app": "configure-a-pod",
				},
				Pods: []client.SimplePod{
					{
						Name:   "nginx-pod",
						Mounts: []client.PodVolume{},
						Labels: map[string]string{
							"app": "configure-a-pod",
						},
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

func (l *Level2) GetClusterStatus(cluster client.Cluster, msg string) string {
	if len(cluster.Deployments) != 1 {
		return "There should be 1 deployment; found: " + string(len(cluster.Deployments))
	}

	deployment := cluster.Deployments[0]
	if len(deployment.Pods) != 1 {
		return "There should be 1 pod; found: " + string(len(deployment.Pods))
	}

	pod := deployment.Pods[0]
	if strings.HasPrefix(pod.Name, "nginx") == false {
		return "Pod was found but named nginx.*; instead a pod named " + pod.Name + " was found"
	}

	if strings.HasPrefix(pod.Containers[0].Name, "nginx") == false {
		return "Container was found but named nginx.*; instead a container named " + pod.Containers[0].Name + " was found"
	}

	if pod.Containers[0].Ports[0] != 80 {
		return "Container port should be 80; found: " + string(pod.Containers[0].Ports[0])
	}

	if pod.Labels["app"] != "configure-a-pod" {
		return "Pod label app should be configure-a-pod; found: " + pod.Labels["app"]
	}

	return "success"
}
