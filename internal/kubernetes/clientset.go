package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"io"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// defaultTailLines sets the logs to get the previous
	// 10000 lines.
	defaultTailLines = int64(10000)
)

//go:generate counterfeiter . Clientset
type Clientset interface {
	PodLogs(string, string, string) (string, error)
	Events(context.Context, string, string, string) ([]v1.Event, error)
}

type clientset struct {
	clientset *kubernetes.Clientset
}

// PodLogs returns logs for a given container in a given pod in a given
// namespace.
func (c *clientset) PodLogs(name, namespace, container string) (string, error) {
	podLogOptions := v1.PodLogOptions{
		Container: container,
		TailLines: &defaultTailLines,
	}

	logs := c.clientset.CoreV1().
		Pods(namespace).
		GetLogs(name, &podLogOptions)

	stream, err := logs.Stream(context.Background())
	if err != nil {
		return "", err
	}
	defer stream.Close()

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, stream)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Events returns events for a given kind, name, and namespace.
func (c *clientset) Events(ctx context.Context, kind, name, namespace string) ([]v1.Event, error) {
	lo := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", name),
		TypeMeta:      metav1.TypeMeta{Kind: kind},
	}

	events, err := c.clientset.CoreV1().Events(namespace).List(ctx, lo)
	if err != nil {
		return nil, err
	}

	return events.Items, nil
}
