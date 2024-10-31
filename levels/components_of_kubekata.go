package levels

import (
	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
	"k8s.io/apimachinery/pkg/labels"
)

/*
	The goal of level1 is to create an nginx deployment
*/

type ComponentsOfKubeKata struct {
	isFinished bool
}

func (l *ComponentsOfKubeKata) GetName() string {
	return "components of kubekata"
}

func (l *ComponentsOfKubeKata) GetDesiredCluster(k8sclient client.Client) client.Cluster {
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

func (l *ComponentsOfKubeKata) GetClusterStatus(cluster client.Cluster, msg string) string {
	return "success"
}

func (l *ComponentsOfKubeKata) SetFinished() {
	l.isFinished = true
}

func (l *ComponentsOfKubeKata) GetIsFinished() bool {
	return l.isFinished
}
