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

// Metrics object to hold metrics data
type PodMetrics struct {
	Metrics         []PodMetricsInfo `json:"metrics"`
	LatestTimestamp time.Time        `json:"latestTimestamp"`
}

type PodMetricsInfo struct {
	Timestamp string `json:"timestamp`
	Value     int64  `json:"value"`
}

// calls heapster to get the metrics of that pod
// https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/podautoscaler/metrics/legacy_metrics_client.go
func getMetricsByPodName(podName string, namespace string) (*PodMetrics, time.Time, error) {

	heapsterHost := "http://localhost:8082" //"http://heapster.kube-system.svc.cluster.local"
	target := fmt.Sprintf("%s/api/v1/model/namespaces/%s/pods/%s/metrics/cpu/usage_rate", heapsterHost, namespace, podName)
	var metrics PodMetrics
	err := getJSON(target, &metrics)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to unmarshal heapster response: %v", err)
	}
	return &metrics, time.Time{}, nil
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

// get containers per pod:
// http://localhost:8082/api/v1/model/debug/allkeys

// [
// "namespace:default/pod:master-identity-identity-3262377241-p8ml8/container:web",
// "namespace:kube-system/pod:weave-cortex-node-exporter-grr6r/container:agent",
// "node:ip-172-22-152-72.eu-west-1.compute.internal/container:kubelet",
// ]

// t1, e := time.Parse(
//         time.RFC3339,
//         "2012-11-01T22:08:41+00:00")

//  curl -k -s -H "Authorization:Bearer $IM_TOKEN" \
// https://master.cfc:8443/kubernetes/api/v1/proxy/\
// namespaces/kube-system/services/heapster/api\
// /v1/model/namespaces/default/pods/ \
// nginx-1720728682-dklzv/metrics/cpu/usage_rate
