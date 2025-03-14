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
	WaitForJobCompletion(ctx context.Context, namespace, jobName string) (err error)
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

// NewDefaultClientInterface creates a new client interface with the default kube client
func NewDefaultClientInterface() (Interface, error) {
	kubeClient, err := Client()
	if err != nil {
		return nil, err
	}
	return NewClientInterface(kubeClient), nil
}

func (c *client) GetClient() kubernetes.Interface {
	return c.kubeClient
}
