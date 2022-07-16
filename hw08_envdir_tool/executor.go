package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	execCmd := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	envs := make([]string, 0)

	for k, v := range env {
		if _, ok := os.LookupEnv(k); ok {
			os.Unsetenv(k)
		}
		if strings.Contains(v, "=") {
			continue
		}
		if v != "" {
			envs = append(envs, fmt.Sprintf("%s=%s", k, v))
		}
	}

	execCmd.Env = append(os.Environ(), envs...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		fmt.Println(err)
		return exitCode
	}

	return execCmd.ProcessState.ExitCode()
}
