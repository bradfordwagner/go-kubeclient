package kube

import (
	"context"
	"k8s.io/client-go/kubernetes"
)

type Interface interface {
	JobInterface
	GetClient() kubernetes.Interface
}

type JobInterface interface {
	DeleteJob(ctx context.Context, namespace, jobName string) (err error)
}

type client struct {
	Interface
	kubeClient kubernetes.Interface
}

func NewClientInterface(kubeClient kubernetes.Interface) Interface {
	return &client{
		kubeClient: kubeClient,
	}
}

func (c *client) GetClient() kubernetes.Interface {
	return c.kubeClient
}
