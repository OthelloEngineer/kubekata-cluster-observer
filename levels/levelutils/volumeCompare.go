package levelutils

import (
	"fmt"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
)

func ComparePodVolumes(expectedVolumes []client.PodVolume, currentVolumes []client.PodVolume) string {
	if len(expectedVolumes) != len(currentVolumes) {
		return fmt.Sprintf("Expected %d volumes, but found %d", len(expectedVolumes), len(currentVolumes))
	}
	for _, v := range expectedVolumes {
		found := false
		msg := "success"
		for _, cv := range currentVolumes {
			msg = comparePodVolume(v, cv)
			if msg == "success" {
				found = true
				break
			}
		}
		if !found {
			return msg
		}
	}
	return "success"
}

func comparePodVolume(expectedVolume client.PodVolume, currentVolume client.PodVolume) string {
	if expectedVolume.Name != currentVolume.Name {
		return fmt.Sprintf("Expected volume name is %s, but found %s", expectedVolume.Name, currentVolume.Name)
	}
	if expectedVolume.PersistentVolumeClaimName != currentVolume.PersistentVolumeClaimName {
		return fmt.Sprintf("Expected volume persistent volume claim name is %s, but found %s", expectedVolume.PersistentVolumeClaimName, currentVolume.PersistentVolumeClaimName)
	}
	return "success"
}
