package kube_test

import (
	"context"
	"time"

	"github.com/bradfordwagner/go-kubeclient/kube"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	authv1 "k8s.io/api/authentication/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	coretesting "k8s.io/client-go/testing"
)

var _ = Describe("CreateServiceAccountToken", func() {
	It("service account does not exist", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, name := "test", "nonexistent"
		timeout, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		token, err := clint.CreateServiceAccountToken(timeout, ns, name, time.Hour)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not found"))
		Expect(token).To(BeEmpty())
	})

	It("creates token successfully", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, name := "test", "test-sa"
		timeout, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err := kubeClient.CoreV1().ServiceAccounts(ns).Create(timeout, &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}, metav1.CreateOptions{})
		Expect(err).To(BeNil())

		expectedToken := "test-token-123"
		kubeClient.PrependReactor("create", "serviceaccounts", func(action coretesting.Action) (handled bool, ret runtime.Object, err error) {
			if action.GetSubresource() == "token" {
				return true, &authv1.TokenRequest{
					Status: authv1.TokenRequestStatus{
						Token: expectedToken,
					},
				}, nil
			}
			return false, nil, nil
		})

		token, err := clint.CreateServiceAccountToken(timeout, ns, name, time.Hour)
		Expect(err).To(BeNil())
		Expect(token).To(Equal(expectedToken))
	})

	It("creates token with specific duration", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, name := "test", "test-sa"
		timeout, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err := kubeClient.CoreV1().ServiceAccounts(ns).Create(timeout, &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}, metav1.CreateOptions{})
		Expect(err).ShouldNot(HaveOccurred())

		expectedToken := "test-token-with-duration"
		var capturedDuration int64
		kubeClient.PrependReactor("create", "serviceaccounts", func(action coretesting.Action) (handled bool, ret runtime.Object, err error) {
			if action.GetSubresource() == "token" {
				createAction := action.(coretesting.CreateAction)
				tokenRequest := createAction.GetObject().(*authv1.TokenRequest)
				capturedDuration = *tokenRequest.Spec.ExpirationSeconds
				return true, &authv1.TokenRequest{
					Status: authv1.TokenRequestStatus{
						Token: expectedToken,
					},
				}, nil
			}
			return false, nil, nil
		})

		duration := 30 * time.Minute
		token, err := clint.CreateServiceAccountToken(timeout, ns, name, duration)
		Expect(err).To(BeNil())
		Expect(token).To(Equal(expectedToken))
		Expect(capturedDuration).To(Equal(int64(duration.Seconds())))
	})

	It("handles zero duration", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, name := "test", "test-sa"
		timeout, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err := kubeClient.CoreV1().ServiceAccounts(ns).Create(timeout, &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}, metav1.CreateOptions{})
		Expect(err).To(BeNil())

		expectedToken := "test-token-zero-duration"
		var capturedDuration int64
		kubeClient.PrependReactor("create", "serviceaccounts", func(action coretesting.Action) (handled bool, ret runtime.Object, err error) {
			if action.GetSubresource() == "token" {
				createAction := action.(coretesting.CreateAction)
				tokenRequest := createAction.GetObject().(*authv1.TokenRequest)
				capturedDuration = *tokenRequest.Spec.ExpirationSeconds
				return true, &authv1.TokenRequest{
					Status: authv1.TokenRequestStatus{
						Token: expectedToken,
					},
				}, nil
			}
			return false, nil, nil
		})

		token, err := clint.CreateServiceAccountToken(timeout, ns, name, 0)
		Expect(err).To(BeNil())
		Expect(token).To(Equal(expectedToken))
		Expect(capturedDuration).To(Equal(int64(0)))
	})

	It("handles context cancellation", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, name := "test", "test-sa"
		ctx, cancel := context.WithCancel(context.Background())

		_, err := kubeClient.CoreV1().ServiceAccounts(ns).Create(ctx, &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}, metav1.CreateOptions{})
		Expect(err).To(BeNil())

		kubeClient.PrependReactor("create", "serviceaccounts", func(action coretesting.Action) (handled bool, ret runtime.Object, err error) {
			if action.GetSubresource() == "token" {
				cancel()
				return true, nil, context.Canceled
			}
			return false, nil, nil
		})

		token, err := clint.CreateServiceAccountToken(ctx, ns, name, time.Hour)
		Expect(err).To(Equal(context.Canceled))
		Expect(token).To(BeEmpty())
	})
})
