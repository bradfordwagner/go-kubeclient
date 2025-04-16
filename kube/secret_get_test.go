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

var _ = Describe("SecretGet", func() {
	It("secret dne", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, sec := "test", "dne"
		timeout, _ := context.WithTimeout(context.Background(), time.Second)

		data, exists, err := clint.SecretGet(timeout, ns, sec)
		Expect(err.Error()).To(ContainSubstring("not found"))
		Expect(exists).To(BeFalse())
		Expect(data).To(BeNil())
	})

	It("finds a secret", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, sec := "test", "exists"
		timeout, _ := context.WithTimeout(context.Background(), time.Second)

		expectedData := map[string]string{
			"hello": "world",
		}
		expectedDataAsBytes := make(map[string][]byte)
		for k, v := range expectedData {
			expectedDataAsBytes[k] = []byte(v)
		}

		kubeClient.CoreV1().Secrets(ns).Create(timeout, &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: sec,
			},
			Data: expectedDataAsBytes,
		}, metav1.CreateOptions{})

		data, exists, err := clint.SecretGet(timeout, ns, sec)
		Expect(err).To(BeNil())
		Expect(exists).To(BeTrue())
		Expect(data).To(Equal(expectedData))
	})
})
