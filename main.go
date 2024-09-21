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

		namespace := getNamespace(configText)
		fmt.Println("assuming namespace: ", namespace)
		kubeclient = client.NewClient("./config", namespace)
		if kubeclient == nil {
			fmt.Print("Error creating client; retrying...")
			kubeclient = client.NewClient("./config", namespace)
		}
		fmt.Println("successfully uploaded config", kubeclient)
		if err != nil {
			w.Write([]byte("Error reading request body"))
			w.WriteHeader(500)
			fmt.Print("Error reading request body")
			return
		}
	})

	http.HandleFunc("/clusterState", func(w http.ResponseWriter, r *http.Request) {
		if kubeclient.Namespace == "" {
			_, err := w.Write([]byte("No namespace set"))
			if err != nil {
				fmt.Print("Error writing response")
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
				fmt.Print("Error writing response")
				return
			}
		}
		num, err := client.DeleteAllResources(*kubeclient)
		if err != nil {
			_, err := w.Write([]byte("Deleted " + string(rune(num)) + " However encountered this error" + err.Error()))
			if err != nil {
				fmt.Print("Error writing response")
				return
			}
		} else {
			_, err := w.Write([]byte("Deleted " + string(rune(num)) + " resources"))
			if err != nil {
				fmt.Print("Error writing response")
				return
			}
		}
	})

	http.HandleFunc("/desired", func(w http.ResponseWriter, r *http.Request) {
		if kubeclient.Namespace == "" {
			_, err := w.Write([]byte("No namespace set"))
			if err != nil {
				fmt.Println("Error writing response")
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

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if kubeclient.Namespace == "" {
			_, err := w.Write([]byte("No namespace set"))
			if err != nil {
				fmt.Print("Error writing response")
				return
			}
		}
		level, err := levelRepository.GetCurrentLevel()
		if err != nil {
			msg := "Error getting current level, likely no level set" + err.Error()
			w.Write([]byte(msg))
			return
		}
		cluster := client.GetAllResources(*kubeclient)
		msg := r.URL.Query().Get("msg")
		diff := level.GetClusterStatus(cluster, msg)
		fmt.Println("Status:", diff)
		if diff == "success" {
			levelRepository.SetCurrentLevel(level.GetID() + 1)
		}
		w.Write([]byte(diff))
	})

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Print("Error starting server," + err.Error())
	}

}

func getNamespace(configText []byte) string {
	nsRegex := `namespace: (.*)`
	compileNsRegex, _ := regexp.Compile(nsRegex)
	namespace := compileNsRegex.FindStringSubmatch(string(configText))
	if len(namespace) < 1 {
		fmt.Println("Could not find namespace using regex 'namespace: (.*)' therefore assuming: namespace: default")
		namespace = []string{"namespace: default"}
	}
	namespace = strings.Split(namespace[0], " ")
	return namespace[1]
}
