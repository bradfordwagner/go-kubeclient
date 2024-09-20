package kube

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/bradfordwagner/go-util/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // this blank import is necessary to load the oidc plugin for client-go: otherwise No Auth Provider found for name "oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func Client(kubeconfig string) (clientset kubernetes.Interface, err error) {
	l := log.Log()
	config, err := config(kubeconfig)
	if err != nil {
		l.With("error", err).Error("failed to create kubernetes config")
		return
	}

	// silence warnings from k8s client-go
	// specifically from gatekeeper
	rest.SetDefaultWarningHandler(rest.NoWarnings{})

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		l.With("error", err).Error("failed to create kubernetes client")
	}

	return
}

func Dynamic(kubeconfig string) (d dynamic.Interface, err error) {
	l := log.Log()

	config, err := config(kubeconfig)
	if err != nil {
		l.With("error", err).Error("failed to create kubernetes config")
		return
	}

	d, err = dynamic.NewForConfig(config)
	if err != nil {
		l.With("error", err).Error("failed to create kubernetes dynamic client")
	}
	return
}

func config(kubeconfig string) (config *rest.Config, err error) {
	// in cluster
	config, err = rest.InClusterConfig()
	if err != nil {
		// kubeconfig / file based
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	return
}

func ScaleDeployment(ctx context.Context, kubeclient kubernetes.Interface, namespace, deploymentName string, replicas int32) (err error) {
	l := log.Log()
	deploymentsClient := kubeclient.AppsV1().Deployments(namespace)
	deployment, err := deploymentsClient.Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		l.With("error", err).Error("failed to get deployment")
		return
	}

	deployment.Spec.Replicas = &replicas
	_, err = deploymentsClient.Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		l.With("error", err).Error("failed to update deployment")
	}
	return
}

func ForceDeletePod(ctx context.Context, kubeclient kubernetes.Interface, namespace, podName string) (err error) {
	l := log.Log()
	podsClient := kubeclient.CoreV1().Pods(namespace)
	err = podsClient.Delete(ctx, podName, metav1.DeleteOptions{})
	if err != nil {
		l.With("error", err).Error("failed to delete pod")
	}
	return
}

func CheckPodEvictedOrContainerStatusUnknown(pod v1.Pod) (isErrored bool) {
	l := log.Log()
	// check for Evicted or ContainerStatusUnknown
	j, err := json.MarshalIndent(pod.Status, "", "  ")
	if err != nil {
		l.With("error", err).Error("failed to marshal pod status")
		return true
	}
	s := string(j)
	if strings.Contains(s, "Evicted") || strings.Contains(s, "ContainerStatusUnknown") {
		return true
	}
	return
}
