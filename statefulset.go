package gokube

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kube) ListStatefulSets(namespace string) ([]string, error) {
	var statefulSetNames []string
	statefulSets, err := k.clientset.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, statefulSet := range statefulSets.Items {
		statefulSetNames = append(statefulSetNames, statefulSet.ObjectMeta.Name)
	}

	return statefulSetNames, nil
}
