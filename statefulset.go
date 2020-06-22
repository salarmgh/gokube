package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func listStatefulSets() []string {
	var statefulSetNames []string
	config, err := GetClientConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := GetClientsetFromConfig(config)
	if err != nil {
		panic(err)
	}
	statefulSets, err := clientset.AppsV1().StatefulSets("app").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, statefulSet := range statefulSets.Items {
		statefulSetNames = append(statefulSetNames, statefulSet.ObjectMeta.Name)
	}

	return statefulSetNames
}
