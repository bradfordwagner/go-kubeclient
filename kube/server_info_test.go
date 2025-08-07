package kube_test

import (
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
		expectedCA := "-----BEGIN CERTIFICATE-----\ntest-ca-certificate-data\n-----END CERTIFICATE-----"

		config := &rest.Config{
			Host: expectedServer,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: []byte(expectedCA),
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

	It("returns raw pem certificate data", func() {
		kubeClient := fake.NewClientset()
		pemCert := `-----BEGIN CERTIFICATE-----
MIIDITCCAgmgAwIBAgIJALnW+rAaFAM5MA0GCSqGSIb3DQEBCwUAMCkxJzAlBgNV
BAMTHmF1dGgwLmF1dGgwLmNvbSBURVNUIENlcnRpZmljYXRlMB4XDTE5MDEwMjE5
MjUyNFoXDTI5MDEwMjE5MjUyNFowKTEnMCUGA1UEAxMeYXV0aDAuYXV0aDAuY29t
IFRFU1QgQ2VydGlmaWNhdGUwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQC7VJTUt9Us8cKBwl+M1M7Bm8Y+gXX7RZUfRUgPm9YjCnYoBjA8cUEKPsyYhiGj
-----END CERTIFICATE-----`

		config := &rest.Config{
			Host: "https://test-cluster.example.com:6443",
			TLSClientConfig: rest.TLSClientConfig{
				CAData: []byte(pemCert),
			},
		}

		client := kube.NewClientInterfaceWithConfig(kubeClient, config)

		_, ca, err := client.GetServerInfo()

		Expect(err).ShouldNot(HaveOccurred())
		Expect(ca).To(Equal(pemCert))
		Expect(ca).To(ContainSubstring("-----BEGIN CERTIFICATE-----"))
		Expect(ca).To(ContainSubstring("-----END CERTIFICATE-----"))
	})
})
