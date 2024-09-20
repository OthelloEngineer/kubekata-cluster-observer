package levels

import (
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
