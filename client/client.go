package client

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type PodVolume struct {
	Name                      string `json:"name"`
	PersistentVolumeClaimName string `json:"persistentVolumeClaimname"`
}

type PersistentVolume struct {
	Name       string `json:"name"`
	Capacity   string `json:"storage"`
	AccessMode string `json:"accessMode"`
}

type PersistentVolumeClaim struct {
	Name             string           `json:"name"`
	Capacity         string           `json:"storage"`
	AccessMod        string           `json:"accessMode"`
	PersistentVolume PersistentVolume `json:"persistentVolume"`
}

type Resource struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
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

type SimplePod struct {
	Name       string      `json:"name"`
	Containers []Container `json:"containers"`
	Mounts     []PodVolume `json:"mounts"`
}

type Deployment struct {
	Name         string            `json:"name"`
	Replicas     int               `json:"replicas"`
	Pods         []SimplePod       `json:"pods"`
	StrategyType string            `json:"strategyType"`
	SelectorMap  map[string]string `json:"selector"`
	Namespace    string            `json:"namespace"`
}

type EndPoint struct {
	Name       string      `json:"name"`
	Containers []SimplePod `json:"pods"`
	Namespace  string      `json:"Namespace"`
}

type Service struct {
	Name        string            `json:"name"`
	Ports       []int32           `json:"ports"`
	SelectorMap map[string]string `json:"selector"`
	Namespace   string            `json:"namespace"`
	Endpoints   []EndPoint        `json:"endpoints"`
}

type Cluster struct {
	Deployments           []Deployment            `json:"deployments"`
	Services              []Service               `json:"services"`
	PersistentVolume      []PersistentVolume      `json:"persistentVolume"`
	PersistentVolumeClaim []PersistentVolumeClaim `json:"persistentVolumeClaim"`
}

const (
	label = "kube=kata"
)

type Client struct {
	Client    *kubernetes.Clientset
	Namespace string
}

func NewClient(configLocation string, namespace string) *Client {
	rawClient, err := clientcmd.BuildConfigFromFlags("", configLocation)
	if err != nil {
		return &Client{}
	}
	clientset, err := kubernetes.NewForConfig(rawClient)
	if err != nil {
		return &Client{}
	}
	// verify that the Client is working
	client := &Client{
		Client:    clientset,
		Namespace: namespace,
	}
	pod, err := clientset.CoreV1().Pods(client.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Print("Likely auth error when getting pods: ", err)
		return &Client{}
	}
	fmt.Print(pod.Items[0].Name)

	return client
}

func GetAllResources(client Client) Cluster {
	listOptions := metav1.ListOptions{
		LabelSelector: label,
	}

	cluster := Cluster{
		Deployments:           []Deployment{},
		Services:              []Service{},
		PersistentVolume:      []PersistentVolume{},
		PersistentVolumeClaim: []PersistentVolumeClaim{},
	}

	pods, err := client.Client.CoreV1().Pods(client.Namespace).List(context.Background(), listOptions)
	fmt.Println("No. pods found: ", len(pods.Items))
	if err != nil {
		return cluster
	}
	var simplePods []SimplePod
	for _, pod := range pods.Items {
		var containers []Container
		for _, container := range pod.Spec.Containers {
			containers = append(containers, Container{
				Name:  container.Name,
				Ports: getPortFromv1Container(container),
				Image: container.Image,
				Requests: Resource{
					CPU:    container.Resources.Requests.Cpu().String(),
					Memory: container.Resources.Requests.Memory().String(),
				},
				Limits: Resource{
					CPU:    container.Resources.Limits.Cpu().String(),
					Memory: container.Resources.Limits.Memory().String(),
				},
				Envs: getEnvsFromV1Container(container),
			})
		}
		newPod := SimplePod{
			Name:       pod.Name,
			Containers: containers,
			Mounts:     getMountsFromV1Pod(pod),
		}
		fmt.Print("Found pod: ", newPod.Name)
		simplePods = append(simplePods, newPod)
	}

	deployments, err := client.Client.AppsV1().Deployments(client.Namespace).List(context.Background(), listOptions)
	if err != nil {
		return cluster
	}
	fmt.Println("No. deployments found: ", len(deployments.Items))
	var deploymentsList []Deployment
	for _, deployment := range deployments.Items {
		deploymentsList = append(deploymentsList, Deployment{
			Name:         deployment.Name,
			Replicas:     int(*deployment.Spec.Replicas),
			Pods:         simplePods,
			StrategyType: string(deployment.Spec.Strategy.Type),
			SelectorMap:  deployment.Spec.Selector.MatchLabels,
			Namespace:    deployment.Namespace,
		})
	}

	cluster.Deployments = deploymentsList

	services, err := client.Client.CoreV1().Services("default").List(context.Background(), listOptions)
	if err != nil {
		return cluster
	}
	fmt.Println("No. services found: ", len(services.Items))
	var servicesList []Service
	for _, service := range services.Items {
		var endpoints []EndPoint
		for _, endpoint := range service.Spec.Ports {
			endpoints = append(endpoints, EndPoint{
				Name:       endpoint.Name,
				Containers: simplePods,
				Namespace:  service.Namespace,
			})
		}
		servicesList = append(servicesList, Service{
			Name:        service.Name,
			Ports:       getPortsFromV1Service(service),
			SelectorMap: service.Spec.Selector,
			Namespace:   service.Namespace,
			Endpoints:   endpoints,
		})
	}

	var persistentVolumes []PersistentVolume
	persistentVolumeList, err := client.Client.CoreV1().PersistentVolumes().List(context.Background(), listOptions)
	if err != nil {
		fmt.Println("Error getting persistent volumes: ", err)
		return cluster
	}

	var persistentVolumeClaims []PersistentVolumeClaim
	persistentVolumeClaimList, err := client.Client.CoreV1().PersistentVolumeClaims(client.Namespace).List(context.Background(), listOptions)
	if err != nil {
		fmt.Println("Error getting persistent volume claims: ", err)
		return cluster
	}
	fmt.Println("No. persistent volume claims found: ", len(persistentVolumeClaimList.Items))
	for _, persistentVolumeClaim := range persistentVolumeClaimList.Items {
		persistentVolumeClaims = append(persistentVolumeClaims, PersistentVolumeClaim{
			Name:      persistentVolumeClaim.Name,
			Capacity:  persistentVolumeClaim.Spec.Resources.Requests.Storage().String(),
			AccessMod: string(persistentVolumeClaim.Spec.AccessModes[0]),
			PersistentVolume: PersistentVolume{
				Name:       persistentVolumeClaim.Spec.VolumeName,
				Capacity:   persistentVolumeClaim.Spec.Resources.Requests.Storage().String(),
				AccessMode: string(persistentVolumeClaim.Spec.AccessModes[0]),
			},
		})
	}

	fmt.Println("No. persistent volumes found: ", len(persistentVolumeList.Items))
	for _, persistentVolume := range persistentVolumeList.Items {
		persistentVolumes = append(persistentVolumes, PersistentVolume{
			Name:       persistentVolume.Name,
			Capacity:   persistentVolume.Spec.Capacity.Storage().String(),
			AccessMode: string(persistentVolume.Spec.AccessModes[0]),
		})
	}

	return cluster

}

func getPortFromv1Container(container v1.Container) []int32 {
	ports := []int32{}
	for _, port := range container.Ports {
		ports = append(ports, port.ContainerPort)
	}
	return ports
}

func getEnvsFromV1Container(container v1.Container) []string {
	envs := []string{}
	for _, env := range container.Env {
		envs = append(envs, env.Name)
	}
	return envs
}

func getMountsFromV1Pod(pod v1.Pod) []PodVolume {
	mounts := []PodVolume{}
	for _, volume := range pod.Spec.Volumes {

		if volume.PersistentVolumeClaim == nil {
			continue
		}

		mounts = append(mounts, PodVolume{
			Name:                      volume.Name,
			PersistentVolumeClaimName: volume.PersistentVolumeClaim.ClaimName,
		})

	}
	return mounts
}

func getPortsFromV1Service(service v1.Service) []int32 {
	ports := []int32{}
	for _, port := range service.Spec.Ports {
		ports = append(ports, port.Port)
	}
	return ports
}

func DeleteAllResources(client Client) (int, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: label,
	}
	numDeletions := 0

	pods, err := client.Client.CoreV1().Pods(client.Namespace).List(context.Background(), listOptions)
	if err != nil {
		fmt.Println("Error getting pods: ", err)
		return 0, err
	}
	fmt.Println("No. pods found: ", len(pods.Items))
	for _, pod := range pods.Items {
		err := client.Client.CoreV1().Pods(client.Namespace).Delete(context.Background(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Println("Error deleting pod: ", err)
			continue
		}
		numDeletions++
	}

	deployments, err := client.Client.AppsV1().Deployments(client.Namespace).List(context.Background(), listOptions)
	if err != nil {
		fmt.Println("Error getting deployments: ", err)
		return numDeletions, err
	}
	fmt.Println("No. deployments found: ", len(deployments.Items))
	for _, deployment := range deployments.Items {
		err := client.Client.AppsV1().Deployments(client.Namespace).Delete(context.Background(), deployment.Name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Println("Error deleting deployment: ", err)
			continue
		}
		numDeletions++
	}

	services, err := client.Client.CoreV1().Services("default").List(context.Background(), listOptions)
	if err != nil {
		fmt.Println("Error getting services: ", err)
		return numDeletions, err
	}
	fmt.Println("No. services found: ", len(services.Items))
	for _, service := range services.Items {
		err := client.Client.CoreV1().Services("default").Delete(context.Background(), service.Name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Println("Error deleting service: ", err)
			continue
		}
		numDeletions++
	}

	persistentVolumeClaims, err := client.Client.CoreV1().PersistentVolumeClaims(client.Namespace).List(context.Background(), listOptions)
	if err != nil {
		fmt.Println("Error getting persistent volume claims: ", err)
		return numDeletions, err
	}
	fmt.Println("No. persistent volume claims found: ", len(persistentVolumeClaims.Items))
	for _, persistentVolumeClaim := range persistentVolumeClaims.Items {
		err := client.Client.CoreV1().PersistentVolumeClaims(client.Namespace).Delete(context.Background(), persistentVolumeClaim.Name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Println("Error deleting persistent volume claim: ", err)
			continue
		}
		numDeletions++
	}

	persistentVolumes, err := client.Client.CoreV1().PersistentVolumes().List(context.Background(), listOptions)
	if err != nil {
		fmt.Println("Error getting persistent volumes: ", err)
		return numDeletions, err
	}
	fmt.Println("No. persistent volumes found: ", len(persistentVolumes.Items))
	for _, persistentVolume := range persistentVolumes.Items {
		err := client.Client.CoreV1().PersistentVolumes().Delete(context.Background(), persistentVolume.Name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Println("Error deleting persistent volume: ", err)
			continue
		}
		numDeletions++
	}

	return numDeletions, nil
}
