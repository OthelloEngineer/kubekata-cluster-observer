package client

type ConfigCluster struct {
	ClusterName string `json:"clusterName"`
	Server      string `json:"server"`
	Name        string `json:"name"`
}

type Context struct {
	ClusterName string `json:"clusterName"`
	Context     string `json:"context"`
	Namespace   string `json:"namespace"`
}

type User struct {
	Name string `json:"name"`
	Key  string `json:"key"`
	Cert string `json:"cert"`
}

type Config struct {
	ClusterNames   []string `json:"clusterNames"`
	Contexts       []string `json:"contexts"`
	CurrentContext string   `json:"currentContext"`
	Namespace      string   `json:"namespace"`
}
