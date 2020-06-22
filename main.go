package main

import (
	"fmt"
)

const debug = false

func main() {
	ExecToPodThroughAPI("./test.sh", "", "accounts-green-0", "app", nil)
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

	ExecToPodThroughAPI(command, containerName, podName, namespace, nil)
}
