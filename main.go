package main

import (
	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
	"github.com/OthelloEngineer/kubekata-cluster-observer/levels"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"os"
)

func main() {

	kubeclient := client.NewClient("./config")
	levelRepository := new(levels.LevelRepository)
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		configText, _ := io.ReadAll(r.Body)
		err := os.WriteFile("config", configText, 0644)
		if err != nil {
			return
		}
	})

	http.HandleFunc("/clusterState", func(w http.ResponseWriter, r *http.Request) {
		cluster := client.GetAllResources(*kubeclient)
		clusterJSON, _ := json.Marshal(cluster)
		w.Write(clusterJSON)
	})

	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		num, err := client.DeleteAllResources(*kubeclient)
		if err != nil {
			_, err := w.Write([]byte("Deleted " + string(rune(num)) + " However encountered this error" + err.Error()))
			if err != nil {
				return
			}
		} else {
			_, err := w.Write([]byte("Deleted " + string(rune(num)) + " resources"))
			if err != nil {
				return
			}
		}
	})

	http.HandleFunc("/desired", func(w http.ResponseWriter, r *http.Request) {
		level, err := levelRepository.GetCurrentLevel()
		if err != nil {
			return
		}
		desiredCluster := level.GetDesiredCluster()
		desiredClusterJSON, _ := json.Marshal(desiredCluster)
		w.Write(desiredClusterJSON)
	})

	http.HandleFunc("/diff", func(w http.ResponseWriter, r *http.Request) {
		level, err := levelRepository.GetCurrentLevel()
		if err != nil {
			return
		}
		cluster := client.GetAllResources(*kubeclient)
		diff := level.GetClusterDiff(cluster)
		w.Write([]byte(diff))
	})

	http.ListenAndServe(":8081", nil)

}
