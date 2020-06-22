package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func getPodsFromStatefulSet(name string) []string {
	var podsList []string
	config, err := GetClientConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := GetClientsetFromConfig(config)
	if err != nil {
		panic(err)
	}
	statefulSet, err := clientset.AppsV1().StatefulSets("app").Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	set := labels.Set(statefulSet.Spec.Selector.MatchLabels)
	pods, err := clientset.CoreV1().Pods("app").List(context.TODO(), metav1.ListOptions{LabelSelector: set.AsSelector().String()})

	for _, pod := range pods.Items {
		podsList = append(podsList, pod.Name)
	}

	return podsList
}
