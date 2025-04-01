package kube

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigmapGet returns the data of a configmap
func (c *client) ConfigmapGet(ctx context.Context, namespace, name string) (data map[string]string, exists bool, err error) {
	get, err := c.kubeClient.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	exists = !errors.IsNotFound(err)
	if !exists {
		return
	}
	data = get.Data

	return
}
