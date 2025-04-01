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

var _ = Describe("ConfigmapDelete", func() {
	It("configmap dne", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, cm := "test", "dne"

		err := clint.ConfigmapDelete(context.Background(), ns, cm)
		Expect(err).To(BeNil())
	})

	It("configmap exists", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, cm := "test", "exists"

		// Create a configmap
		kubeClient.CoreV1().ConfigMaps(ns).Create(context.Background(), &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: cm,
			},
		}, metav1.CreateOptions{})

		err := clint.ConfigmapDelete(context.Background(), ns, cm)
		Expect(err).To(BeNil())
	})
})
