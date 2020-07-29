package gokube

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (k *Kube) GetPodsFromStatefulSet(statefulset string, namespace string) ([]string, error) {
	var podsList []string
	statefulSet, err := k.clientset.AppsV1().StatefulSets(namespace).Get(statefulset, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	set := labels.Set(statefulSet.Spec.Selector.MatchLabels)
	pods, err := k.clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: set.AsSelector().String()})
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {
		podsList = append(podsList, pod.Name)
	}
	return podsList, nil
}
