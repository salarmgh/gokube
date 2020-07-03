package gokube

import (
	"k8s.io/client-go/kubernetes"
)

type Kube struct {
	clientset *kubernetes.Clientset
}
