package gokube

import (
	batchv1 "k8s.io/api/batch/v1"
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

func (k *Kube) CreateJob(name string, namespace string, image string, command []string, args []string, labels map[string]string, env map[string]string, volumes []core.Volume, mounts []core.VolumeMount) error {
	var envs []core.EnvVar
	for k, v := range env {
		envs = append(envs, core.EnvVar{Name: k, Value: v})
	}
	jobConfig := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels:    labels,
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:         name,
							Image:        image,
							Command:      command,
							Args:         args,
							Env:          envs,
							VolumeMounts: mounts,
						},
					},
					RestartPolicy: "Never",
					Volumes:       volumes,
				},
			},
		},
	}
	_, err := k.clientset.BatchV1().Jobs(jobConfig.Namespace).Create(jobConfig)
	if err != nil {
		return err
	}
	return nil
}

func (k *Kube) DeleteJob(name string, namespace string) error {
	deleteOptions := metav1.DeleteOptions{}
	err := k.clientset.BatchV1().Jobs(namespace).Delete(name, &deleteOptions)
	if err != nil {
		return err
	}
	return nil
}

func (k *Kube) JobStatus(name string, namespace string) (string, error) {
	getOptions := metav1.GetOptions{}
	job, err := k.clientset.BatchV1().Jobs(namespace).Get(name, getOptions)
	if err != nil {
		return "", err
	}
	if job.Status.Active > 0 {
		return "Running", nil
	} else {
		if job.Status.Succeeded > 0 {
			return "Completed", nil
		}
	}
	return "Failed", nil
}
