package kube_test

import (
	"context"
	"github.com/bradfordwagner/go-kubeclient/kube"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("SecretDelete", func() {
	It("secret dne", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, sec := "test", "dne"

		err := clint.SecretDelete(context.Background(), ns, sec)
		Expect(err).To(BeNil())
	})

	It("secret exists", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, sec := "test", "exists"

		// Create a configmap
		kubeClient.CoreV1().Secrets(ns).Create(context.Background(), &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: sec,
			},
		}, metav1.CreateOptions{})

		err := clint.SecretDelete(context.Background(), ns, sec)
		Expect(err).To(BeNil())
	})
})
