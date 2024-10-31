package levels

import "github.com/OthelloEngineer/kubekata-cluster-observer/client"

type ExposingToTheWorld struct {
	isFinished bool
}

func (l *ExposingToTheWorld) GetName() string {
	return "exposing to the world"
}

func (l *ExposingToTheWorld) SetFinished() {
	l.isFinished = true
}

func (l *ExposingToTheWorld) GetIsFinished() bool {
	return l.isFinished
}

func (l *ExposingToTheWorld) GetDesiredCluster() client.Cluster {
	deployment := expectedDeployment()

	service := getExpectedService()
	service.Type = "NodePort" // Quite naive, maybe a better goal is to be considered?

	return client.Cluster{
		Deployments:           []client.Deployment{deployment},
		Services:              []client.Service{service},
		PersistentVolume:      *new([]client.PersistentVolume),
		PersistentVolumeClaim: *new([]client.PersistentVolumeClaim),
	}
}
