package kube

import (
	"context"
	authv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// CreateToken creates a token for a service account
// returns the token string and an error if the token could not be created
func (c *client) CreateServiceAccountToken(ctx context.Context, namespace, name string, duration time.Duration) (token string, err error) {
	durationSeconds := int64(duration.Seconds())

	tokenRequest := &authv1.TokenRequest{
		Spec: authv1.TokenRequestSpec{
			ExpirationSeconds: &durationSeconds,
		},
	}

	result, err := c.kubeClient.CoreV1().ServiceAccounts(namespace).CreateToken(ctx, name, tokenRequest, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	return result.Status.Token, nil
}
