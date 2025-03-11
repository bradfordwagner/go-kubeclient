package kube

import (
	"context"
	"fmt"
	"github.com/bradfordwagner/go-util/bwutil"
	"github.com/bradfordwagner/go-util/log"
	"golang.org/x/sync/errgroup"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

// DeleteJob deletes a job synchronously
// returns an error if we could not delete the job
func (c *client) DeleteJob(ctx context.Context, namespace, jobName string) (err error) {
	l := log.Log().With("action", "delete", "namespace", namespace, "job", jobName)
	_, err = c.kubeClient.BatchV1().Jobs(namespace).Get(ctx, jobName, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		l.Info("job not found")
		return nil
	}

	// use err group to bubble up errors
	errgroup, ctx := errgroup.WithContext(ctx)

	// watch job
	errgroup.Go(func() (err error) {
		watcher, err := c.kubeClient.BatchV1().Jobs(namespace).Watch(ctx, metav1.ListOptions{
			FieldSelector:  "metadata.name=" + jobName,
			TimeoutSeconds: bwutil.Pointer[int64](365 * 24 * 60 * 60), // 1 year
		})
		if err != nil {
			return
		}
		defer watcher.Stop()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case event, open := <-watcher.ResultChan():
				if event.Type == watch.Deleted || !open {
					return
				}
			}
		}
	})

	// delete job
	errgroup.Go(func() error {
		propagation := metav1.DeletePropagationForeground
		return c.kubeClient.BatchV1().Jobs(namespace).Delete(ctx, jobName, metav1.DeleteOptions{
			// DeletePropagationForeground - will force propagation of deletion job -> pods
			PropagationPolicy: &propagation,
		})
	})

	// collect results
	err = errgroup.Wait()
	if err != nil {
		l.With("error", err).Error("failed to delete job")
	} else {
		l.Info("job deleted")
	}
	return
}
