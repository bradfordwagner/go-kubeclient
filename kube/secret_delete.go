package kube

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretDelete deletes a secret from the specified namespace and name.
func (c *client) SecretDelete(ctx context.Context, namespace, name string) (err error) {
	err = c.kubeClient.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		err = nil
	}
	return
}
