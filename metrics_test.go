package main

import (
	"fmt"
	"testing"
)

func TestGetPods(t *testing.T) {

	labels :=
		map[string]string{
			"k8s-app": "heapster",
		}

	pods, err := getPods(labels, "kube-system")
	if err != nil {
		t.Errorf("Error getting pods %v", err)
	}
	for _, pod := range pods.Items {
		if pod.Status.Phase == "Running" {
			fmt.Printf("%v\n", pod.Name)
		}
	}
	if len(pods.Items) == 0 {
		t.Errorf("Error getting heapster %v", err)
	}
}

func TestGetMetrics(t *testing.T) {
	namespace := "default"
	labels :=
		map[string]string{
			"component": "controller",
		}

	pods, err := getPods(labels, namespace)
	if err != nil {
		t.Errorf("Error getting pods %v", err)
	}
	for _, pod := range pods.Items {
		metrics, ts, err := getMetricsByPodName(pod.Name, namespace)
		if err != nil {
			t.Errorf("Error getting metrics %v", err)
		}
		fmt.Printf("%v -- %v\n", ts, metrics)

	}
	if len(pods.Items) == 0 {
		t.Errorf("Error getting heapster %v", err)
	}
}
