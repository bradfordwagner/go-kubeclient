package kube

import (
	"context"
	"fmt"
	"github.com/bradfordwagner/go-util/bwutil"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// StatefulsetEvent wraps a watch.Event and the StatefulSet object.
type StatefulsetEvent struct {
	EventType   watch.EventType
	Statefulset *appsv1.StatefulSet
	RawEvent    watch.Event
}

// StatefulSetWatcher provides access to StatefulSet watch events.
type StatefulSetWatcher interface {
	ResultChan() <-chan StatefulsetEvent
	Stop()
}

type statefulSetWatcher struct {
	events chan StatefulsetEvent
	stop   func()
}

func (s *statefulSetWatcher) ResultChan() <-chan StatefulsetEvent {
	return s.events
}

func (s *statefulSetWatcher) Stop() {
	s.stop()
}

// WatchStatefulset watches a StatefulSet by name in the given namespace.
func (c *client) WatchStatefulset(ctx context.Context, namespace, name string) (StatefulSetWatcher, error) {
	watcher, err := c.kubeClient.AppsV1().StatefulSets(namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector:  "metadata.name=" + name,
		TimeoutSeconds: bwutil.Pointer[int64](365 * 24 * 60 * 60), // 1 year
	})
	if err != nil {
		return nil, fmt.Errorf("failed to watch statefulset: %w", err)
	}

	events := make(chan StatefulsetEvent)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer close(events)
		defer watcher.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case event, open := <-watcher.ResultChan():
				if !open {
					return
				}
				ss, _ := event.Object.(*appsv1.StatefulSet)
				events <- StatefulsetEvent{
					EventType:   event.Type,
					Statefulset: ss,
					RawEvent:    event,
				}
			}
		}
	}()

	return &statefulSetWatcher{
		events: events,
		stop:   cancel,
	}, nil
}
