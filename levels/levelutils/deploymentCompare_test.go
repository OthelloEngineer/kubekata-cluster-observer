package levelutils_test

import (
	"fmt"
	"testing"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
	"github.com/OthelloEngineer/kubekata-cluster-observer/levels/levelutils"
)

func TestIdenticalClusters(t *testing.T) {
	simplePod1 := client.NewSimplePod(
		"pod1",
		[]client.Container{
			client.NewContainer(
				"container1",
				"image1",
				[]int32{8080},
				client.NewResource("1", "1Gi"),
				client.NewResource("1", "1Gi"),
				[]string{"env1=env1", "env2=env2"},
				[]string{"mount1"},
			),
		},
		[]client.PodVolume{
			client.NewPodVolume("volume1", "pvc1"),
		},
		map[string]string{
			"app": "app1",
		},
	)

	deploy1 := client.NewDeployment(
		"deploy1",
		"rolling",
		"namespace1",
		1,
		[]client.SimplePod{simplePod1},
		map[string]string{
			"app": "app1",
		},
	)

	svc1 := client.NewService(
		"svc1",
		"namespace1",
		[]int32{8080},
		map[string]string{
			"app": "app1",
		},
		[]client.EndPoint{
			client.NewEndPoint("ep1", "namespace1", simplePod1),
		},
		"NodePort",
	)

	pv1 := client.NewPersistentVolume("pv1", "1Gi", "rwo")

	pvc1 := client.NewPersistentVolumeClaim("pvc1", "1Gi", "rwo", pv1)

	currentCluster := client.Cluster{
		Deployments:           []client.Deployment{deploy1},
		Services:              []client.Service{svc1},
		PersistentVolume:      []client.PersistentVolume{pv1},
		PersistentVolumeClaim: []client.PersistentVolumeClaim{pvc1},
	}

	expectedCluster := client.Cluster{
		Deployments:           []client.Deployment{deploy1},
		Services:              []client.Service{svc1},
		PersistentVolume:      []client.PersistentVolume{pv1},
		PersistentVolumeClaim: []client.PersistentVolumeClaim{pvc1},
	}

	result := levelutils.CompareDeployments(currentCluster.Deployments[0], expectedCluster.Deployments[0])
	if result != "success" {
		t.Errorf("Expected success, got %s", result)
	}
}

func TestDifferentDeploy(t *testing.T) {
	simplePod1 := client.NewSimplePod(
		"pod1",
		[]client.Container{
			client.NewContainer(
				"container1",
				"image1",
				[]int32{8080},
				client.NewResource("1", "1Gi"),
				client.NewResource("1", "1Gi"),
				[]string{"env1=env1", "env2=env2"},
				[]string{"mount1"},
			),
		},
		[]client.PodVolume{
			client.NewPodVolume("volume1", "pvc1"),
		},
		map[string]string{
			"app": "app1",
		},
	)

	deploy1 := client.NewDeployment(
		"deploy1",
		"rolling",
		"namespace1",
		1,
		[]client.SimplePod{simplePod1},
		map[string]string{
			"app": "app1",
		},
	)

	deploy2 := client.NewDeployment(
		"deploy2",
		"rolling",
		"namespace1",
		1,
		[]client.SimplePod{simplePod1},
		map[string]string{
			"app": "app1",
		},
	)

	svc1 := client.NewService(
		"svc1",
		"namespace1",
		[]int32{8080},
		map[string]string{
			"app": "app1",
		},
		[]client.EndPoint{
			client.NewEndPoint("ep1", "namespace1", simplePod1),
		},
		"NodePort",
	)

	pv1 := client.NewPersistentVolume("pv1", "1Gi", "rwo")

	pvc1 := client.NewPersistentVolumeClaim("pvc1", "1Gi", "rwo", pv1)

	currentCluster := client.Cluster{
		Deployments:           []client.Deployment{deploy1},
		Services:              []client.Service{svc1},
		PersistentVolume:      []client.PersistentVolume{pv1},
		PersistentVolumeClaim: []client.PersistentVolumeClaim{pvc1},
	}

	expectedCluster := client.Cluster{
		Deployments:           []client.Deployment{deploy2},
		Services:              []client.Service{svc1},
		PersistentVolume:      []client.PersistentVolume{pv1},
		PersistentVolumeClaim: []client.PersistentVolumeClaim{pvc1},
	}

	result := levelutils.CompareDeployments(currentCluster.Deployments[0], expectedCluster.Deployments[0])
	if result != "Deployment name should be deploy2; found: deploy1" {
		t.Errorf("Expected Deployment name should be deploy2; found: deploy1, got %s", result)
	}
}

func TestCompareImagesAndPortOfDeployments(t *testing.T) {
	simplePod1 := client.NewSimplePod(
		"pod1",
		[]client.Container{
			client.NewContainer(
				"container1",
				"image1",
				[]int32{8080},
				client.NewResource("1", "1Gi"),
				client.NewResource("1", "1Gi"),
				[]string{"env1=env1", "env2=env2"},
				[]string{"mount1"},
			),
		},
		[]client.PodVolume{
			client.NewPodVolume("volume1", "pvc1"),
		},
		map[string]string{
			"app": "app1",
		},
	)

	deploy1 := client.NewDeployment(
		"deploy1",
		"rolling",
		"namespace1",
		1,
		[]client.SimplePod{simplePod1},
		*new(map[string]string),
	)

	simplePod2 := client.NewSimplePod(
		"pod1",
		[]client.Container{
			client.NewContainer(
				"container1",
				"image2",
				[]int32{8081},
				client.NewResource("1", "1Gi"),
				client.NewResource("1", "1Gi"),
				[]string{"env1=env1", "env2=env2"},
				[]string{"mount1"},
			),
		},
		*new([]client.PodVolume),
		map[string]string{
			"app": "app1",
		},
	)

	deploy2 := client.NewDeployment(
		"deploy1",
		"rolling",
		"namespace1",
		1,
		[]client.SimplePod{simplePod2},
		*new(map[string]string),
	)

	msg := levelutils.CompareImagesAndPortOfDeployments([]client.Deployment{deploy1}, []client.Deployment{deploy2})
	if msg != "Container image image2 not found" {
		t.Errorf("Expected Container image should be image2; found: image1, got %s", msg)
	}
}

func TestCompareImagesAndPortOfDeploymentsDifferentPorts(t *testing.T) {
	simplePod1 := client.NewSimplePod(
		"pod1",
		[]client.Container{
			client.NewContainer(
				"container1",
				"image1",
				[]int32{8080},
				client.NewResource("1", "1Gi"),
				client.NewResource("1", "1Gi"),
				[]string{"env1=env1", "env2=env2"},
				[]string{"mount1"},
			),
		},
		[]client.PodVolume{
			client.NewPodVolume("volume1", "pvc1"),
		},
		map[string]string{
			"app": "app1",
		},
	)

	deploy1 := client.NewDeployment(
		"deploy1",
		"rolling",
		"namespace1",
		1,
		[]client.SimplePod{simplePod1},
		*new(map[string]string),
	)

	simplePod2 := client.NewSimplePod(
		"pod1",
		[]client.Container{
			client.NewContainer(
				"container1",
				"image1",
				[]int32{8081},
				client.NewResource("1", "1Gi"),
				client.NewResource("1", "1Gi"),
				[]string{"env1=env1", "env2=env2"},
				[]string{"mount1"},
			),
		},
		*new([]client.PodVolume),
		map[string]string{
			"app": "app1",
		},
	)

	deploy2 := client.NewDeployment(
		"deploy1",
		"rolling",
		"namespace1",
		1,
		[]client.SimplePod{simplePod2},
		*new(map[string]string),
	)

	msg := levelutils.CompareImagesAndPortOfDeployments([]client.Deployment{deploy1}, []client.Deployment{deploy2})
	expectedMsg := fmt.Sprintf("Container image %s should have port %d; found: %d", simplePod2.Containers[0].Image, simplePod2.Containers[0].Ports[0], simplePod1.Containers[0].Ports[0])
	if msg != expectedMsg {
		t.Errorf("Expected %s, got %s", expectedMsg, msg)
	}
}
