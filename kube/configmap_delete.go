package kube

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigmapDelete deletes a configmap
func (c *client) ConfigmapDelete(ctx context.Context, namespace, name string) (err error) {
	err = c.kubeClient.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		err = nil
	}
	return
}
