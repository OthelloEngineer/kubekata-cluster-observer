package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
	"github.com/OthelloEngineer/kubekata-cluster-observer/levels"
	"k8s.io/apimachinery/pkg/util/json"
)

func main() {
	kubeclient := new(client.Client)
	kubeclient.Namespace = ""
	levelRepository := levels.NewLevelRepository()
	levelRepository.SetCurrentLevel(1)
	currentLevel, _ := levelRepository.GetCurrentLevel()
	fmt.Println("Starting server with current level", currentLevel)
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		configText, _ := io.ReadAll(r.Body)
		err := os.WriteFile("config", configText, 0644)
		if err != nil {
			fmt.Print("Error writing config file")
			return
		}
		nsRegex := `namespace: (.*)`
		compileNsRegex, _ := regexp.Compile(nsRegex)
		namespace := compileNsRegex.FindStringSubmatch(string(configText))[0]
		namespace = strings.Split(namespace, " ")[1]
		fmt.Println("Namespace found in config", namespace)
		if namespace == "" {
			fmt.Print("Namespace not found in config")
			return
		}
		kubeclient = client.NewClient("./config", namespace)
		fmt.Println("successfully uploaded config", kubeclient)
		if err != nil {
			return
		}
	})

	http.HandleFunc("/clusterState", func(w http.ResponseWriter, r *http.Request) {
		if kubeclient.Namespace == "" {
			_, err := w.Write([]byte("No namespace set"))
			if err != nil {
				return
			}
		}
		cluster := client.GetAllResources(*kubeclient)
		clusterJSON, _ := json.Marshal(cluster)
		w.Write(clusterJSON)
	})

	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		if kubeclient.Namespace == "" {
			_, err := w.Write([]byte("No namespace set"))
			if err != nil {
				return
			}
		}
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
		if kubeclient.Namespace == "" {
			_, err := w.Write([]byte("No namespace set"))
			if err != nil {
				return
			}
		}
		level, err := levelRepository.GetCurrentLevel()
		if err != nil {
			_, _ = w.Write([]byte("Error getting current level, likely no level set"))
			return
		}
		desiredCluster := level.GetDesiredCluster()
		desiredClusterJSON, _ := json.Marshal(desiredCluster)
		w.Write(desiredClusterJSON)
	})

	http.HandleFunc("/diff", func(w http.ResponseWriter, r *http.Request) {
		if kubeclient.Namespace == "" {
			_, err := w.Write([]byte("No namespace set"))
			if err != nil {
				return
			}
		}
		level, err := levelRepository.GetCurrentLevel()
		if err != nil {
			w.Write([]byte("Error getting current level, likely no level set"))
			return
		}
		cluster := client.GetAllResources(*kubeclient)
		diff := level.GetClusterDiff(cluster)
		w.Write([]byte(diff))
	})

	http.ListenAndServe(":8080", nil)

}
