package levels

import "github.com/OthelloEngineer/kubekata-cluster-observer/client"

type ScalingTheApp struct {
	isFinished bool
}

func (l *ScalingTheApp) GetName() string {
	return "scaling the app"
}

func (l *ScalingTheApp) SetFinished() {
	l.isFinished = true
}

func (l *ScalingTheApp) GetIsFinished() bool {
	return l.isFinished
}

func (l *ScalingTheApp) GetDesiredCluster() client.Cluster {
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

func (l *ScalingTheApp) GetClusterStatus(cluster client.Cluster, msg string, k8sclient client.Client) string {
	if cluster.Deployments[0].Name != "hello-go" {
		return "There should be a deployment named hello-go"
	}

	if len(cluster.Deployments) != 1 {
		return "There should be 1 deployment; found: " + string(rune(len(cluster.Deployments)))
	}

	if len(cluster.Deployments[0].Pods) != 3 {
		return "There should be 3 pods; found: " + string(rune(len(cluster.Deployments[0].Pods)))
	}

	return "success"
}
