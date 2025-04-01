package kube_test

import (
	"context"
	"github.com/bradfordwagner/go-kubeclient/kube"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"time"
)

var _ = Describe("ConfigmapGet", func() {
	It("configmap dne", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, cm := "test", "dne"
		timeout, _ := context.WithTimeout(context.Background(), time.Second)

		data, exists, err := clint.ConfigmapGet(timeout, ns, cm)
		Expect(err.Error()).To(ContainSubstring("not found"))
		Expect(exists).To(BeFalse())
		Expect(data).To(BeNil())
	})

	It("finds a configmap", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, cm := "test", "dne"
		timeout, _ := context.WithTimeout(context.Background(), time.Second)

		expectedData := map[string]string{
			"hello": "world",
		}
		kubeClient.CoreV1().ConfigMaps(ns).Create(timeout, &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: cm,
			},
			Data: expectedData,
		}, metav1.CreateOptions{})

		data, exists, err := clint.ConfigmapGet(timeout, ns, cm)
		Expect(err).To(BeNil())
		Expect(exists).To(BeTrue())
		Expect(data).To(Equal(expectedData))
	})
})
