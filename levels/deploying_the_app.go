package levels

import (
	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
	"github.com/OthelloEngineer/kubekata-cluster-observer/levels/levelutils"
)

type DeployingTheApp struct {
	isFinished bool
}

func (l *DeployingTheApp) GetName() string {
	return "simple deployment"
}

func expectedDeployment() client.Deployment {
	deployment := client.NewDeployment(
		"hello-go",
		"app",
		"hello-go",
		1,
		[]client.SimplePod{
			client.NewSimplePod("hello-go", []client.Container{
				client.NewContainer("hello-go", "othelloengineer/hello-go:1.0.0", []int32{8080}, *new(client.Resource), *new(client.Resource), []string{}, []string{}),
			},
				[]client.PodVolume{}, map[string]string{}),
		},
		map[string]string{"app": "hello-go"},
	)
	return deployment
}

func (l *DeployingTheApp) GetDesiredCluster() client.Cluster {
	return client.Cluster{
		Deployments:           []client.Deployment{expectedDeployment()},
		Services:              *new([]client.Service),
		PersistentVolume:      *new([]client.PersistentVolume),
		PersistentVolumeClaim: *new([]client.PersistentVolumeClaim),
	}
}

func (l *DeployingTheApp) GetClusterStatus(cluster client.Cluster, msg string) string {
	result := levelutils.CompareImagesAndPortOfDeployments([]client.Deployment{expectedDeployment()}, cluster.Deployments)
	return result
}

func (l *DeployingTheApp) SetFinished() {
	l.isFinished = true
}

func (l *DeployingTheApp) GetIsFinished() bool {
	return l.isFinished
}
