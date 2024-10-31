package levels

import (
	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
	"github.com/OthelloEngineer/kubekata-cluster-observer/levels/levelutils"
)

type dns_and_services struct {
	isFinished bool
}

func (l *dns_and_services) GetName() string {
	return "dns and services"
}

func (l *dns_and_services) SetFinished() {
	l.isFinished = true
}

func (l *dns_and_services) GetDesiredCluster(k8sclient client.Client) client.Cluster {
	dep := expectedDeployment()
	svc := getExpectedService()
	return client.Cluster{
		Deployments:           []client.Deployment{dep},
		Services:              []client.Service{svc},
		PersistentVolume:      []client.PersistentVolume{},
		PersistentVolumeClaim: []client.PersistentVolumeClaim{},
	}
}

func getExpectedService() client.Service {
	service := client.NewService(
		"hello-go",
		"",
		[]int32{8080},
		map[string]string{
			"app": "hello-go",
		},
		[]client.EndPoint{
			client.NewEndPoint("hello-go", "", expectedDeployment().Pods[0]),
		},
		"NodePort",
	)
	return service
}

func (l *dns_and_services) GetClusterStatus(cluster client.Cluster, msg string) string {
	if len(cluster.Services) != 1 {
		return "There should be 1 service; found: " + string(len(cluster.Services))
	}

	statusMsg := levelutils.CompareServices(cluster.Services[0], getExpectedService(), getExpectedService().Endpoints)
	if statusMsg != "success" {
		return statusMsg
	}

	statusMsg = levelutils.CompareDeployments(cluster.Deployments[0], expectedDeployment())
	if statusMsg != "success" {
		return statusMsg
	}

	return statusMsg
}

func (l *dns_and_services) GetIsFinished() bool {
	return l.isFinished
}
