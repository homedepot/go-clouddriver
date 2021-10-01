package kubernetes

import (
	"bytes"
	"context"
	"io"

	v1 "k8s.io/api/core/v1"
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
