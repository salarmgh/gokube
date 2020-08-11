package gokube

import (
	core "k8s.io/api/core/v1"
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

func (k *Kube) CreatePod(name string, namespace string, image string, command []string) error {
	pod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"digiops": "true",
				"logger":  "true",
			},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            name,
					Image:           image,
					ImagePullPolicy: core.PullIfNotPresent,
					Command:         command,
				},
			},
		},
	}
	pod, err := k.clientset.CoreV1().Pods(pod.Namespace).Create(pod)
	if err != nil {
		return err
	}
	return nil
}
