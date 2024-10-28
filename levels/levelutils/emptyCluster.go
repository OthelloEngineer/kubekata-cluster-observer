package levelutils

import "github.com/OthelloEngineer/kubekata-cluster-observer/client"

func GetEmptyCluster() client.Cluster {
	deployments := []client.Deployment{}
	services := []client.Service{}
	persistentVolumes := []client.PersistentVolume{}
	persistentVolumeClaims := []client.PersistentVolumeClaim{}

	return client.Cluster{
		Deployments:           deployments,
		Services:              services,
		PersistentVolume:      persistentVolumes,
		PersistentVolumeClaim: persistentVolumeClaims,
	}
}
