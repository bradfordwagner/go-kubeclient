package kube_test

import (
	"encoding/base64"
	"github.com/bradfordwagner/go-kubeclient/kube"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

var _ = Describe("GetServerInfo", func() {
	It("returns server and ca when config is available", func() {
		kubeClient := fake.NewClientset()
		expectedServer := "https://test-cluster.example.com:6443"
		expectedCABytes := []byte("test-ca-certificate-data")
		expectedCA := base64.StdEncoding.EncodeToString(expectedCABytes)

		config := &rest.Config{
			Host: expectedServer,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: expectedCABytes,
			},
		}

		client := kube.NewClientInterfaceWithConfig(kubeClient, config)

		server, ca, err := client.GetServerInfo()

		Expect(err).ShouldNot(HaveOccurred())
		Expect(server).To(Equal(expectedServer))
		Expect(ca).To(Equal(expectedCA))
	})

	It("returns empty values when config has empty data", func() {
		kubeClient := fake.NewClientset()
		config := &rest.Config{
			Host: "",
			TLSClientConfig: rest.TLSClientConfig{
				CAData: nil,
			},
		}

		client := kube.NewClientInterfaceWithConfig(kubeClient, config)

		server, ca, err := client.GetServerInfo()

		Expect(err).ShouldNot(HaveOccurred())
		Expect(server).To(BeEmpty())
		Expect(ca).To(BeEmpty())
	})

	It("returns error when config is not available", func() {
		kubeClient := fake.NewClientset()
		client := kube.NewClientInterface(kubeClient)

		server, ca, err := client.GetServerInfo()

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("rest config is not available"))
		Expect(server).To(BeEmpty())
		Expect(ca).To(BeEmpty())
	})

	It("handles config with only server but no ca data", func() {
		kubeClient := fake.NewClientset()
		expectedServer := "https://minimal-cluster.example.com:6443"

		config := &rest.Config{
			Host: expectedServer,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: nil,
			},
		}

		client := kube.NewClientInterfaceWithConfig(kubeClient, config)

		server, ca, err := client.GetServerInfo()

		Expect(err).ShouldNot(HaveOccurred())
		Expect(server).To(Equal(expectedServer))
		Expect(ca).To(BeEmpty())
	})

	It("handles config with only ca data but no server", func() {
		kubeClient := fake.NewClientset()
		expectedCABytes := []byte("standalone-ca-certificate-data")
		expectedCA := base64.StdEncoding.EncodeToString(expectedCABytes)

		config := &rest.Config{
			Host: "",
			TLSClientConfig: rest.TLSClientConfig{
				CAData: expectedCABytes,
			},
		}

		client := kube.NewClientInterfaceWithConfig(kubeClient, config)

		server, ca, err := client.GetServerInfo()

		Expect(err).ShouldNot(HaveOccurred())
		Expect(server).To(BeEmpty())
		Expect(ca).To(Equal(expectedCA))
	})

	It("handles large ca certificate data", func() {
		kubeClient := fake.NewClientset()
		expectedServer := "https://prod-cluster.example.com:6443"
		// Create a large CA data to test handling of bigger certificates
		expectedCABytes := make([]byte, 4096)
		for i := range expectedCABytes {
			expectedCABytes[i] = byte(i % 256)
		}
		expectedCA := base64.StdEncoding.EncodeToString(expectedCABytes)

		config := &rest.Config{
			Host: expectedServer,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: expectedCABytes,
			},
		}

		client := kube.NewClientInterfaceWithConfig(kubeClient, config)

		server, ca, err := client.GetServerInfo()

		Expect(err).ShouldNot(HaveOccurred())
		Expect(server).To(Equal(expectedServer))
		Expect(ca).To(Equal(expectedCA))
		Expect(len(ca)).To(BeNumerically(">", 0))
	})
})
