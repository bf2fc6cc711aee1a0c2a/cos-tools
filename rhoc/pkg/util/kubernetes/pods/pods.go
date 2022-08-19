package pods

import (
	"bufio"
	"context"
	"errors"
	"io"
	"regexp"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var empty []byte
var re = regexp.MustCompile(ansi)

func ListContainers(ctx context.Context, client kubernetes.Interface, namespace string, name string) ([]string, error) {
	result, err := client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	containers := make([]string, 0)

	for _, container := range result.Spec.Containers {
		containers = append(containers, container.Name)
	}

	return containers, nil
}

func Logs(ctx context.Context, client kubernetes.Interface, namespace string, name string, container string, out io.Writer) error {
	logOptions := corev1.PodLogOptions{
		Follow:    false,
		Container: container,
	}

	byteReader, err := client.CoreV1().Pods(namespace).GetLogs(name, &logOptions).Stream(ctx)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(byteReader)
	for {
		data, err := reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			return err
		}
		if err != nil {
			break
		}

		data = re.ReplaceAll(data, empty)

		if _, err = out.Write(data); err != nil {
			break
		}
	}

	return nil
}
