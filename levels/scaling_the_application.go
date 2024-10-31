package levels

import "github.com/OthelloEngineer/kubekata-cluster-observer/client"

type ScalingTheApplication struct {
	isFinished bool
}

func (l *ScalingTheApplication) GetName() string {
	return "scaling the application"
}

func (l *ScalingTheApplication) SetFinished() {
	l.isFinished = true
}

func (l *ScalingTheApplication) GetIsFinished() bool {
	return l.isFinished
}

func (l *ScalingTheApplication) GetDesiredCluster() client.Cluster {
	pod := expectedDeployment().Pods[0]
	pods := []client.SimplePod{pod, pod, pod}
	deployment := expectedDeployment()

	deployment.Pods = pods

	return client.Cluster{
		Deployments:           []client.Deployment{deployment},
		Services:              *new([]client.Service),
		PersistentVolume:      *new([]client.PersistentVolume),
		PersistentVolumeClaim: *new([]client.PersistentVolumeClaim),
	}
}

func (l *ScalingTheApplication) GetClusterStatus(cluster client.Cluster, msg string) string {

	if cluster.Deployments[0].Name != "hello-go" {
		return "There should be a deployment named hello-go"
	}

	if len(cluster.Deployments) != 1 {
		return "There should be 1 deployment; found: " + string(len(cluster.Deployments))
	}

	if len(cluster.Deployments[0].Pods) != 3 {
		return "There should be 3 pods; found: " + string(len(cluster.Deployments[0].Pods))
	}

	return "success"
}
