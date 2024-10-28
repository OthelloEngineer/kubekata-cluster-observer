package client

type PodVolume struct {
	Name                      string `json:"name"`
	PersistentVolumeClaimName string `json:"persistentVolumeClaimname"`
}

func NewPodVolume(name, pvcName string) PodVolume {
	return PodVolume{
		Name:                      name,
		PersistentVolumeClaimName: pvcName,
	}
}

type PersistentVolume struct {
	Name       string `json:"name"`
	Capacity   string `json:"storage"`
	AccessMode string `json:"accessMode"`
}

func NewPersistentVolume(name, capacity, accessMode string) PersistentVolume {
	return PersistentVolume{
		Name:       name,
		Capacity:   capacity,
		AccessMode: accessMode,
	}
}

type PersistentVolumeClaim struct {
	Name             string           `json:"name"`
	Capacity         string           `json:"storage"`
	AccessMod        string           `json:"accessMode"`
	PersistentVolume PersistentVolume `json:"persistentVolume"`
}

func NewPersistentVolumeClaim(name, capacity, accessMode string, pv PersistentVolume) PersistentVolumeClaim {
	return PersistentVolumeClaim{
		Name:             name,
		Capacity:         capacity,
		AccessMod:        accessMode,
		PersistentVolume: pv,
	}
}

type Resource struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

func NewResource(cpu, memory string) Resource {
	return Resource{
		CPU:    cpu,
		Memory: memory,
	}
}

type Container struct {
	Name     string   `json:"name"`
	Ports    []int32  `json:"ports"`
	Image    string   `json:"image"`
	Requests Resource `json:"requests"`
	Limits   Resource `json:"limits"`
	Envs     []string `json:"envs"`
	Mount    []string `json:"mounts"`
}

func NewContainer(name, image string, ports []int32, requests, limits Resource, envs, mounts []string) Container {
	return Container{
		Name:     name,
		Ports:    ports,
		Image:    image,
		Requests: requests,
		Limits:   limits,
		Envs:     envs,
		Mount:    mounts,
	}
}

type SimplePod struct {
	Name       string            `json:"name"`
	Containers []Container       `json:"containers"`
	Mounts     []PodVolume       `json:"mounts"`
	Labels     map[string]string `json:"labels"`
}

func NewSimplePod(name string, containers []Container, mounts []PodVolume, labels map[string]string) SimplePod {
	return SimplePod{
		Name:       name,
		Containers: containers,
		Mounts:     mounts,
		Labels:     labels,
	}
}

type Deployment struct {
	Name         string            `json:"name"`
	Replicas     int               `json:"replicas"`
	Pods         []SimplePod       `json:"pods"`
	StrategyType string            `json:"strategyType"`
	SelectorMap  map[string]string `json:"selector"`
	Namespace    string            `json:"namespace"`
}

func NewDeployment(name, strategyType, namespace string, replicas int, pods []SimplePod, selector map[string]string) Deployment {
	return Deployment{
		Name:         name,
		Replicas:     replicas,
		Pods:         pods,
		StrategyType: strategyType,
		SelectorMap:  selector,
		Namespace:    namespace,
	}
}

type EndPoint struct {
	Name      string    `json:"name"`
	Pod       SimplePod `json:"pod"`
	Namespace string    `json:"Namespace"`
}

func NewEndPoint(name, namespace string, pod SimplePod) EndPoint {
	return EndPoint{
		Name:      name,
		Pod:       pod,
		Namespace: namespace,
	}
}

type Service struct {
	Name        string            `json:"name"`
	Ports       []int32           `json:"ports"`
	SelectorMap map[string]string `json:"selector"`
	Type        string            `json:"type"`
	Namespace   string            `json:"namespace"`
	Endpoints   []EndPoint        `json:"endpoints"`
}

func NewService(name, namespace string, ports []int32, selector map[string]string, endpoints []EndPoint, svcType string) Service {
	return Service{
		Name:        name,
		Ports:       ports,
		SelectorMap: selector,
		Namespace:   namespace,
		Endpoints:   endpoints,
		Type:        svcType,
	}
}

type Cluster struct {
	Deployments           []Deployment            `json:"deployments"`
	Services              []Service               `json:"services"`
	PersistentVolume      []PersistentVolume      `json:"persistentVolume"`
	PersistentVolumeClaim []PersistentVolumeClaim `json:"persistentVolumeClaim"`
}

func NewCluster(deployments []Deployment, services []Service, pv []PersistentVolume, pvc []PersistentVolumeClaim) Cluster {
	return Cluster{
		Deployments:           deployments,
		Services:              services,
		PersistentVolume:      pv,
		PersistentVolumeClaim: pvc,
	}
}
