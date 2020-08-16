package gokube

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (k *Kube) GetServices(namespace string) (*[]string, error) {
	var services []string
	listOptions := metav1.ListOptions{}
	svcs, err := k.clientset.CoreV1().Services(namespace).List(listOptions)
	if err != nil {
		return nil, err
	}
	for _, svc := range svcs.Items {
		services = append(services, svc.Name)
	}
	return &services, nil
}

func (k *Kube) GetService(name string, namespace string, selector map[string]string) (*apiv1.Service, error) {
	set := labels.Set(selector)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	svcs, err := k.clientset.CoreV1().Services(namespace).List(listOptions)
	if err != nil {
		return nil, err
	}
	return &svcs.Items[0], nil
}

func (k *Kube) GetDeployment(name string, namespace string, selector map[string]string) (*appsv1.Deployment, error) {
	set := labels.Set(selector)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	deployments, err := k.clientset.AppsV1().Deployments(namespace).List(listOptions)
	if err != nil {
		return nil, err
	}
	return &deployments.Items[0], nil
}

func (k *Kube) GetActiveEnv(deployment string, namespace string) (string, error) {
	svc, err := k.GetService(fmt.Sprintf("%s-active", deployment), namespace, map[string]string{"app.kubernetes.io/instance": deployment})
	if err != nil {
		return "", err
	}
	return svc.Spec.Selector["app.env"], nil
}

func (k *Kube) GetActiveEnvDeployment(deployment string, namespace string) (*appsv1.Deployment, error) {
	env, err := k.GetActiveEnv(deployment, namespace)
	if err != nil {
		return nil, err
	}
	dep, err := k.GetDeployment(deployment, namespace, map[string]string{"app.env": env, "app.kubernetes.io/instance": deployment})
	return dep, nil
}
