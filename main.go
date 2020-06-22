package main

import (
	"fmt"
)

const debug = false

func main() {
	pods := getPodsFromStatefulSet("accounts-blue")
	fmt.Printf("%v\n", pods)
}

func testExec() {
	var namespace, containerName, podName, command string
	fmt.Print("Enter namespace: ")
	fmt.Scanln(&namespace)
	fmt.Print("Enter name of the pod: ")
	fmt.Scanln(&podName)
	fmt.Print("Enter name of the container [leave empty if there is only one container]: ")
	fmt.Scanln(&containerName)
	fmt.Print("Enter the commmand to execute: ")
	fmt.Scanln(&command)

	output, stderr, err := ExecToPodThroughAPI(command, containerName, podName, namespace, nil)

	if len(stderr) != 0 {
		fmt.Println("STDERR:", stderr)
	}
	if err != nil {
		fmt.Printf("Error occured while `exec`ing to the Pod %q, namespace %q, command %q. Error: %+v\n", podName, namespace, command, err)
	} else {
		fmt.Println("Output:")
		fmt.Println(output)
	}
}
