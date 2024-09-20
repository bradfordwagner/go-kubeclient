package kube

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Client", func() {

	It("should create a kubernetes client", func() {
		clientset, err := Client()
		Expect(err).ToNot(HaveOccurred())
		Expect(clientset).ToNot(BeNil())
	})

	It("will list all pods in all namespaces", func() {
		clientset, err := Client()
		Expect(err).ToNot(HaveOccurred())
		Expect(clientset).ToNot(BeNil())

		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(pods).ToNot(BeNil())
		Expect(len(pods.Items)).ToNot(BeZero())
	})
})
