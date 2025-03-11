package kube_test

import (
	"context"
	"github.com/bradfordwagner/go-kubeclient/kube"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/sync/errgroup"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake"
	k8stest "k8s.io/client-go/testing"
	"time"
)

var _ = Describe("JobWaitForCompletion", func() {
	It("watches a job to completion", func() {
		kubeClient := fake.NewClientset()
		ns, job := "default", "job"
		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		watcher := watch.NewFake()
		kubeClient.PrependWatchReactor("jobs", k8stest.DefaultWatchReactor(watcher, nil))

		// create the job
		jobObject, err := kubeClient.BatchV1().Jobs(ns).Create(
			ctx,
			&batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name: job,
				},
			},
			metav1.CreateOptions{},
		)
		Expect(err).NotTo(HaveOccurred())

		// invoke the test
		errgroup, ctx := errgroup.WithContext(ctx)
		errgroup.Go(func() (err error) {
			return kube.NewClientInterface(kubeClient).WaitForJobCompletion(ctx, ns, job)
		})

		// emit completion event to watcher
		errgroup.Go(func() (err error) {
			jobObject.Status.Conditions = append(jobObject.Status.Conditions, batchv1.JobCondition{
				Type:   batchv1.JobComplete,
				Status: corev1.ConditionTrue,
			})
			watcher.Modify(jobObject)
			return
		})

		// check results
		Expect(errgroup.Wait()).To(Succeed())
	})
	It("job fails error is bubbled up", func() {
		kubeClient := fake.NewClientset()
		ns, job := "default", "job"
		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		watcher := watch.NewFake()
		kubeClient.PrependWatchReactor("jobs", k8stest.DefaultWatchReactor(watcher, nil))

		// create the job
		jobObject, err := kubeClient.BatchV1().Jobs(ns).Create(
			ctx,
			&batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name: job,
				},
			},
			metav1.CreateOptions{},
		)
		Expect(err).NotTo(HaveOccurred())

		// invoke the test
		errgroup, ctx := errgroup.WithContext(ctx)
		errgroup.Go(func() (err error) {
			return kube.NewClientInterface(kubeClient).WaitForJobCompletion(ctx, ns, job)
		})

		// emit completion event to watcher
		errgroup.Go(func() (err error) {
			jobObject.Status.Conditions = append(jobObject.Status.Conditions, batchv1.JobCondition{
				Type:   batchv1.JobFailed,
				Status: corev1.ConditionTrue,
			})
			watcher.Modify(jobObject)
			return
		})

		// check results
		err = errgroup.Wait()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("job failed: "))
	})
	It("receives a watch error", func() {
		kubeClient := fake.NewClientset()
		ns, job := "default", "job"
		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		watcher := watch.NewFake()
		kubeClient.PrependWatchReactor("jobs", k8stest.DefaultWatchReactor(watcher, nil))

		// create the job
		jobObject, err := kubeClient.BatchV1().Jobs(ns).Create(
			ctx,
			&batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name: job,
				},
			},
			metav1.CreateOptions{},
		)
		Expect(err).NotTo(HaveOccurred())

		// invoke the test
		errgroup, ctx := errgroup.WithContext(ctx)
		errgroup.Go(func() (err error) {
			return kube.NewClientInterface(kubeClient).WaitForJobCompletion(ctx, ns, job)
		})

		// emit completion event to watcher
		errgroup.Go(func() (err error) {
			// error the watch
			watcher.Action(watch.Error, jobObject)
			return
		})

		// check results
		err = errgroup.Wait()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("error watching job"))
	})

	It("times out", func() {
		kubeClient := fake.NewClientset()
		ns, job := "default", "job"
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)

		watcher := watch.NewFake()
		kubeClient.PrependWatchReactor("jobs", k8stest.DefaultWatchReactor(watcher, nil))

		err := kube.NewClientInterface(kubeClient).WaitForJobCompletion(ctx, ns, job)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("context deadline exceeded"))
	})

})
