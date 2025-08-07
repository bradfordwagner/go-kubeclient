package kube

import (
	"context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
)

type Interface interface {
	ConfigmapInterface
	JobInterface
	SecretInterface
	ServiceAccountInterface
	StatefulSetInterface
	GetClient() kubernetes.Interface
	GetServerInfo() (server string, ca string, err error)
}

type JobInterface interface {
	DeleteJob(ctx context.Context, namespace, jobName string) (err error)
	WaitForJobCompletion(ctx context.Context, namespace, jobName string) (err error)
}

type ConfigmapInterface interface {
	ConfigmapGet(ctx context.Context, namespace, name string) (data map[string]string, exists bool, err error)
	ConfigmapDelete(ctx context.Context, namespace, name string) (err error)
}

type SecretInterface interface {
	SecretGet(ctx context.Context, namespace, name string) (data map[string]string, exists bool, err error)
	SecretDelete(ctx context.Context, namespace, name string) (err error)
}

type StatefulSetInterface interface {
	WatchStatefulset(ctx context.Context, namespace, name string) (watcher StatefulSetWatcher, err error)
}

type ServiceAccountInterface interface {
	CreateServiceAccountToken(ctx context.Context, namespace, name string, duration time.Duration) (token string, err error)
}

type client struct {
	Interface
	kubeClient kubernetes.Interface
	config     *rest.Config
}

func NewClientInterface(kubeClient kubernetes.Interface) Interface {
	return &client{
		kubeClient: kubeClient,
	}
}

func NewClientInterfaceWithConfig(kubeClient kubernetes.Interface, config *rest.Config) Interface {
	return &client{
		kubeClient: kubeClient,
		config:     config,
	}
}

// NewDefaultClientInterface creates a new client interface with the default kube client
func NewDefaultClientInterface() (clint Interface, kubeClient kubernetes.Interface, err error) {
	kubeconfig := defaultKubeConfig()
	config, err := config(kubeconfig)
	if err != nil {
		return
	}

	kubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	clint = NewClientInterfaceWithConfig(kubeClient, config)
	return
}

func (c *client) GetClient() kubernetes.Interface {
	return c.kubeClient
}
