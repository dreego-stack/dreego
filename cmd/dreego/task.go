package main

import (
	"fmt"
	"os"
	"os/exec"
)

func cmdTask(args []string) {
	taskBin, err := exec.LookPath("task")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: 'task' executable not found on PATH.\n")
		fmt.Fprintf(os.Stderr, "  dreego task forwards to the Task runner (https://taskfile.dev).\n")
		fmt.Fprintf(os.Stderr, "  Install it from https://taskfile.dev/installation/ and retry.\n")
		os.Exit(1)
	}

	if len(args) == 0 {
		args = []string{"--list"}
	}

	c := exec.Command(taskBin, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			os.Exit(exit.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "task error: %v\n", err)
		os.Exit(1)
	}
}
