package kube

import (
	"context"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// WaitForJobCompletion waits for a job to complete
func WaitForJobCompletion(ctx context.Context, client kubernetes.Interface, namespace, jobName string) error {
	watcher, err := client.BatchV1().Jobs(namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + jobName,
	})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	for {
		select {
		case event := <-watcher.ResultChan():
			if event.Type == watch.Error {
				return fmt.Errorf("error watching job")
			}
			if job, ok := event.Object.(*batchv1.Job); ok {
				for _, condition := range job.Status.Conditions {
					if condition.Type == batchv1.JobComplete && condition.Status == corev1.ConditionTrue {
						// Job completed successfully
						return nil
					}
					if condition.Type == batchv1.JobFailed && condition.Status == corev1.ConditionTrue {
						return fmt.Errorf("job failed: %v", condition.Reason)
					}
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
