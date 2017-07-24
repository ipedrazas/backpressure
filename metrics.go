package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// Get the pod names that match the service selector
func getPods(lbls map[string]string, namespace string) (*v1.PodList, error) {

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	selector := labels.SelectorFromSet(labels.Set(lbls)).String()
	return clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: selector,
	})

}

// calls heapster to get the metrics of that pod
// TODO: add time window
// https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/podautoscaler/metrics/legacy_metrics_client.go
func getMetricsByPodName(podName string, namespace string) (map[string]interface{}, time.Time, error) {

	heapsterHost := "http://localhost:8082" //"http://heapster.kube-system.svc.cluster.local"
	target := fmt.Sprintf("%s/api/v1/model/namespaces/%s/pods/%s/metrics/cpu/usage_rate", heapsterHost, namespace, podName)
	var metrics map[string]interface{}
	err := getJSON(target, &metrics)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to unmarshal heapster response: %v", err)
	}
	return metrics, time.Time{}, nil
}

func getJSON(url string, target interface{}) error {
	myClient := &http.Client{Timeout: 10 * time.Second}
	fmt.Println(url)
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	fmt.Println(r.Body)
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
