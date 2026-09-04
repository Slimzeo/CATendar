//go:build !windows

package agent

import "os/exec"

func configureProcess(command *exec.Cmd) {}
