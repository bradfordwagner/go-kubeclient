package kube

import (
	"context"
	"k8s.io/client-go/kubernetes"
)

type Interface interface {
	JobInterface
}

type JobInterface interface {
	DeleteJob(ctx context.Context, namespace, jobName string) error
}

type client struct {
	Interface
	kubeClient kubernetes.Interface
}

func NewClient(kubeClient kubernetes.Interface) Interface {
	return &client{
		kubeClient: kubeClient,
	}
}
