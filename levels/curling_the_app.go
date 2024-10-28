package levels

import (
	"fmt"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
)

type CurlingTheApp struct {
	isFinished bool
}

func (l *CurlingTheApp) GetName() string {
	return "curling the app"
}

func (l *CurlingTheApp) GetDesiredCluster() client.Cluster {

	return client.Cluster{
		Deployments:           []client.Deployment{expectedDeployment()},
		Services:              *new([]client.Service),
		PersistentVolume:      *new([]client.PersistentVolume),
		PersistentVolumeClaim: *new([]client.PersistentVolumeClaim),
	}
}

func (l *CurlingTheApp) GetClusterStatus(cluster client.Cluster, msg string) string {
	fmt.Println("msg: ", msg)
	if msg == "COOL MESSAGE!!!" {
		return "success"
	}
	return "Expecting the right message :)"
}

func (l *CurlingTheApp) SetFinished() {
	l.isFinished = true
}
