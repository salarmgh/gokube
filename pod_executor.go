package gokube

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
)

func ExecToPodThroughAPI(command, containerName, podName, namespace string, stdin io.Reader) {
	es, _ := elasticsearch.NewDefaultClient()
	config, err := GetClientConfig()
	if err != nil {
		log.Println(err)
	}

	clientset, err := GetClientsetFromConfig(config)
	if err != nil {
		log.Println(err)
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := core_v1.AddToScheme(scheme); err != nil {
		log.Println(err)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&core_v1.PodExecOptions{
		Command:   strings.Fields(command),
		Container: containerName,
		Stdin:     stdin != nil,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		log.Println(err)
	}

	//var stdout, stderr bytes.Buffer
	stdout := new(Buffer)
	stderr := new(bytes.Buffer)
	ch := make(chan bool)
	go func(stdin *bytes.Buffer, stdout *Buffer, stderr *bytes.Buffer, ch chan bool) {
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: stdout,
			Stderr: stderr,
			Tty:    false,
		})
		if err != nil {
			log.Println(err)
			ch <- false
		}
		ch <- true
	}(nil, stdout, stderr, ch)
	buf := bufio.NewReader(stdout)
	isFinished := false
	for !isFinished {
		select {
		case _, ok := <-ch:
			if ok {
				isFinished = true
			}
		default:
			line, _, _ := buf.ReadLine()
			if len(line) == 0 {
				time.Sleep(time.Second)
				break
			}
			var b strings.Builder
			b.WriteString(`{"message" : "`)
			b.WriteString(string(line))
			b.WriteString(`"}`)
			req := esapi.IndexRequest{
				Index: "test",
				Body:  strings.NewReader(b.String()),
			}
			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()
			time.Sleep(time.Second)
		}
	}
}
