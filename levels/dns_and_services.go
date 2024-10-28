package levels

import "github.com/OthelloEngineer/kubekata-cluster-observer/client"

type dns_and_services struct {
	isFinished bool
}

func (l *dns_and_services) GetName() string {
	return "dns and services"
}

func (l *dns_and_services) SetFinished() {
	l.isFinished = true
}

func (l *dns_and_services) GetDesiredCluster() client.Cluster {
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
		[]client.EndPoint{},
		"NodePort",
	)
	return service
}

func (l *dns_and_services) GetClusterStatus(cluster client.Cluster, msg string) string {
	result := compareDeploymentAndService(expectedDeployment(), getExpectedService(), cluster.Deployments, cluster.Services)
	return result
}
