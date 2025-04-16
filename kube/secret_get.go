package kube

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretGet retrieves a secret from the specified namespace and name.
func (c *client) SecretGet(ctx context.Context, namespace, name string) (data map[string]string, exists bool, err error) {
	get, err := c.kubeClient.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	exists = !errors.IsNotFound(err)
	if !exists {
		return
	}
	// convert to string
	data = make(map[string]string)
	for k, bytes := range get.Data {
		data[k] = string(bytes)
	}

	return
}
