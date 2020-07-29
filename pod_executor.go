package gokube

import (
	"bufio"
	"io"
	"time"

	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
)

type Result struct {
	Message string
	Type    string
	Error   error
}

func (k *Kube) ExecToPod(command []string, containerName string, podName string,
	namespace string, resultChan chan Result) {
	signal := make(chan bool)
	defer close(resultChan)
	defer close(signal)
	req := k.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := core_v1.AddToScheme(scheme); err != nil {
		resultChan <- Result{"", "", err}
		return
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&core_v1.PodExecOptions{
		Command:   command,
		Container: containerName,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	config, err := k.getClientConfig()
	if err != nil {
		resultChan <- Result{"", "", err}
		return
	}

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		resultChan <- Result{"", "", err}
		return
	}

	var stdout, stderr Buffer
	go func(stdout *Buffer, stderr *Buffer, signal chan bool) {
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: stdout,
			Stderr: stderr,
			Tty:    false,
		})
		if err != nil {
			resultChan <- Result{"", "", err}
			return
		}
		signal <- true
	}(&stdout, &stderr, signal)

	stdoutReader := bufio.NewReader(&stdout)
	stderrReader := bufio.NewReader(&stderr)

	isFinished := false
	for !isFinished {
		select {
		case _, ok := <-signal:
			if ok {
				isFinished = true
				stdoutResponse := getMessage(stdoutReader, "stdout")
				if stdoutResponse.Error != io.EOF {
					resultChan <- stdoutResponse
				}
				stderrResponse := getMessage(stderrReader, "stderr")
				if stderrResponse.Error != io.EOF {
					resultChan <- stderrResponse
				}
			}
		default:
			stdoutResponse := getMessage(stdoutReader, "stdout")
			if stdoutResponse.Error != io.EOF {
				resultChan <- stdoutResponse
			}
			stderrResponse := getMessage(stderrReader, "stderr")
			if stderrResponse.Error != io.EOF {
				resultChan <- stderrResponse
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func getMessage(stdReader *bufio.Reader, msgType string) Result {
	stdResponse, err := flushLog(stdReader)
	if err != nil {
		return Result{"", "", err}
	}
	return Result{Message: stdResponse, Type: msgType, Error: nil}
}

func flushLog(reader *bufio.Reader) (string, error) {
	p := make([]byte, 1024*100)
	n, err := reader.Read(p)
	if n != 0 && err == nil {
		return string(p[:n]), nil
	} else if err != nil && err == io.EOF && n != 0 {
		return string(p[:n]), nil
	} else {
		return "", err
	}
}
