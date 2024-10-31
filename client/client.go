package client

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	label = "kube=kata"
)

type Client struct {
	Client    *kubernetes.Clientset
	Namespace string
}

func NewClientFromConfig(configLocation string, namespace string) *Client {
	rawClient, err := clientcmd.BuildConfigFromFlags("", configLocation)
	if err != nil {
		fmt.Println("Error building config from flags: ", err)
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
	if len(pod.Items) == 0 {
		fmt.Print("No pods found in namespace: ", namespace)
		return client
	}
	fmt.Print(pod.Items[0].Name)

	return client
}

func ClientFromServiceAccount(namespace string) *Client {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("Error building in cluster config: ", err)
		return &Client{}
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error building client from config: ", err)
		return &Client{}
	}

	pod := client.CoreV1().Pods(namespace)
	if pod == nil {
		fmt.Println("Error getting pods: ", err)
		return &Client{
			Client:    client,
			Namespace: namespace,
		}
	}
	podList, err := pod.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error getting pods: ", err)
		return &Client{
			Client:    client,
			Namespace: namespace,
		}
	}

	fmt.Println("No. pods found: ", len(podList.Items))

	return &Client{
		Client:    client,
		Namespace: namespace,
	}

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

	pods, err := client.Client.CoreV1().Pods(client.Namespace).List(context.Background(), *new(metav1.ListOptions))
	fmt.Println("No. pods found: ", len(pods.Items))
	if err != nil {
		fmt.Println("Error getting pods: ", err)
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
			Labels:     pod.Labels,
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
		deploymentPods := []SimplePod{}
		for _, pod := range simplePods {
			labels := pod.Labels
			if labels == nil {
				labels = make(map[string]string)
			}
			fmt.Println("Pod labels: ", labels)
			for key, value := range deployment.Spec.Selector.MatchLabels {
				fmt.Println("deployment - Key: ", key, "Value: ", value)
				fmt.Println("pod - Key: ", key, "Value: ", labels[key])
				if labels[key] == value {
					deploymentPods = append(deploymentPods, pod)
				}
			}
		}
		deploymentsList = append(deploymentsList, Deployment{
			Name:         deployment.Name,
			Replicas:     int(*deployment.Spec.Replicas),
			Pods:         deploymentPods,
			StrategyType: string(deployment.Spec.Strategy.Type),
			SelectorMap:  deployment.Spec.Selector.MatchLabels,
			Namespace:    deployment.Namespace,
		})

	}

	cluster.Deployments = deploymentsList

	//strayPods := findStrayPods(deploymentsList, simplePods) // this is a crude way of tracking unmanaged pods
	//if len(strayPods) > 0 {
	//	fmt.Println("Stray pods found: ", len(strayPods))
	//	strayDeployment := Deployment{
	//		Name:         "Stray Pods",
	//		Replicas:     len(strayPods),
	//		Pods:         strayPods,
	//		StrategyType: "None",
	//		SelectorMap:  map[string]string{},
	//		Namespace:    client.Namespace,
	//	}
	//	cluster.Deployments = append(cluster.Deployments, strayDeployment)
	//}

	endpointSlices, err := client.Client.DiscoveryV1().EndpointSlices(client.Namespace).List(context.Background(), listOptions)
	endpoints := []EndPoint{}
	if err != nil {
		fmt.Println("Error getting endpoint slices: ", err)
	}

	fmt.Println("No. endpoint slices found: ", len(endpointSlices.Items))

	for _, endpointSlice := range endpointSlices.Items {
		for _, endpoint := range endpointSlice.Endpoints {
			podName := strings.ToLower(endpoint.TargetRef.Name)
			TrimmedPodName := strings.Replace(podName, "pod/", "", 1)
			endpoints = append(endpoints, NewEndPoint(endpointSlice.Name, endpointSlice.Namespace, getPodByName(simplePods, TrimmedPodName)))
		}
	}

	services, err := client.Client.CoreV1().Services("default").List(context.Background(), listOptions)
	if err != nil {
		return cluster
	}
	fmt.Println("No. services found: ", len(services.Items))

	var servicesList []Service
	for _, service := range services.Items {
		endpoints = getEndpointsOfService(service, endpoints)
		servicesList = append(servicesList, Service{
			Name:        service.Name,
			Ports:       getPortsFromV1Service(service),
			SelectorMap: service.Spec.Selector,
			Namespace:   service.Namespace,
			Endpoints:   endpoints,
		})
	}

	cluster.Services = servicesList

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

	cluster.PersistentVolumeClaim = persistentVolumeClaims

	var persistentVolumes []PersistentVolume
	persistentVolumeList, err := client.Client.CoreV1().PersistentVolumes().List(context.Background(), listOptions)

	if err != nil {
		fmt.Println("Error getting persistent volumes: ", err)
		return cluster
	}
	fmt.Println("No. persistent volumes found: ", len(persistentVolumeList.Items))
	for _, persistentVolume := range persistentVolumeList.Items {
		persistentVolumes = append(persistentVolumes, PersistentVolume{
			Name:       persistentVolume.Name,
			Capacity:   persistentVolume.Spec.Capacity.Storage().String(),
			AccessMode: string(persistentVolume.Spec.AccessModes[0]),
		})
	}

	cluster.PersistentVolume = persistentVolumes

	return cluster

}

func getPodByName(pods []SimplePod, name string) SimplePod {
	for _, pod := range pods {
		if pod.Name == name {
			return pod
		}
	}
	return SimplePod{}
}

func getEndpointsOfService(service v1.Service, endpoints []EndPoint) []EndPoint {
	serviceEndpoints := []EndPoint{}
	for _, endpoint := range endpoints {
		EndPointServiceIdentifier := strings.Split(endpoint.Name, "-")[0]
		if service.Name == EndPointServiceIdentifier {
			serviceEndpoints = append(serviceEndpoints, endpoint)
		}
	}
	return serviceEndpoints
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

func findStrayPods(deployments []Deployment, pods []SimplePod) []SimplePod {
	var strayPods []SimplePod
	for _, pod := range pods {
		found := false
		for _, deployment := range deployments {
			for _, deploymentPod := range deployment.Pods {
				if pod.Name == deploymentPod.Name {
					found = true
					break
				}
			}
		}
		if !found {
			strayPods = append(strayPods, pod)
		}
	}
	return strayPods
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
