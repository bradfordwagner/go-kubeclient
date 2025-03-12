package kube_test

import (
	"context"
	"github.com/bradfordwagner/go-kubeclient/kube"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/sync/errgroup"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake"
	k8stest "k8s.io/client-go/testing"
	"time"
)

var _ = Describe("JobDelete", func() {
	It("will delete a job", func() {
		kubeClient := fake.NewClientset()
		ns, job := "default", "job"
		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		watcher := watch.NewFake()
		kubeClient.PrependWatchReactor("jobs", k8stest.DefaultWatchReactor(watcher, nil))

		// create the job
		_, err := kubeClient.BatchV1().Jobs(ns).Create(
			ctx,
			&batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name: job,
				},
			},
			metav1.CreateOptions{},
		)
		Expect(err).NotTo(HaveOccurred())

		eg, ctx := errgroup.WithContext(ctx)
		eg.Go(func() error {
			// run client to delete it
			return kube.NewClientInterface(kubeClient).DeleteJob(ctx, ns, job)
		})

		// emit deleted event to watcher
		eg.Go(func() (err error) {
			watcher.Delete(&batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name: job,
				},
				Spec:   batchv1.JobSpec{},
				Status: batchv1.JobStatus{},
			})
			return
		})

		// success!
		Expect(eg.Wait()).To(Succeed())
	})

	It("will timeout while deleting", func() {
		kubeClient := fake.NewClientset()
		ns, job := "default", "job"
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)

		watcher := watch.NewFake()
		kubeClient.PrependWatchReactor("jobs", k8stest.DefaultWatchReactor(watcher, nil))

		// create the job
		_, err := kubeClient.BatchV1().Jobs(ns).Create(
			ctx,
			&batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name: job,
				},
			},
			metav1.CreateOptions{},
		)
		Expect(err).NotTo(HaveOccurred())

		eg, ctx := errgroup.WithContext(ctx)
		eg.Go(func() error {
			// run client to delete it
			return kube.NewClientInterface(kubeClient).DeleteJob(ctx, ns, job)
		})

		// do not emit any events to watcher, to force timeout

		// failure -- timeout!!
		err = eg.Wait()
		Expect(err).ShouldNot(Succeed())
		Expect(err.Error()).To(ContainSubstring("context deadline exceeded"))
	})

	It("job dne returns no error", func() {
		kubeClient := fake.NewClientset()
		clint := kube.NewClientInterface(kubeClient)
		ns, job := "default", "job"
		ctx := context.Background()

		Expect(clint.DeleteJob(ctx, ns, job)).ToNot(HaveOccurred())
	})
})
