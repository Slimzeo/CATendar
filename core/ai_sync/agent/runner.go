package agent

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type ProgressSink func()

type Runner struct{}

func (Runner) Run(ctx context.Context, adapter Adapter, request RunRequest, progress ProgressSink) error {
	invocation := adapter.Build(request)
	command := exec.CommandContext(ctx, invocation.Path, invocation.Args...)
	command.Dir = request.WorkDir
	command.Env = safeEnvironment(request.Environment)
	command.Stdin = strings.NewReader(invocation.Stdin)
	configureProcess(command)

	stdout, err := command.StdoutPipe()
	if err != nil {
		return fmt.Errorf("capture %s output: %w", adapter.Provider(), err)
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		return fmt.Errorf("capture %s errors: %w", adapter.Provider(), err)
	}
	if err := command.Start(); err != nil {
		return fmt.Errorf("start %s CLI: %w", adapter.Provider(), err)
	}

	var waitGroup sync.WaitGroup
	var errorTail string
	var errorMu sync.Mutex
	var progressOnce sync.Once
	waitGroup.Add(2)
	go func() {
		defer waitGroup.Done()
		drainOutput(stdout, func() {
			if progress != nil {
				progressOnce.Do(progress)
			}
		})
	}()
	go func() {
		defer waitGroup.Done()
		tail := readTail(stderr, 8*1024)
		errorMu.Lock()
		errorTail = tail
		errorMu.Unlock()
	}()

	waitErr := command.Wait()
	waitGroup.Wait()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if waitErr != nil {
		errorMu.Lock()
		detail := lastNonEmptyLine(errorTail)
		errorMu.Unlock()
		if len(detail) > 400 {
			detail = detail[:400]
		}
		if detail != "" {
			return fmt.Errorf("%s CLI failed: %s", adapter.Provider(), detail)
		}
		return fmt.Errorf("%s CLI failed: %w", adapter.Provider(), waitErr)
	}
	return nil
}

func drainOutput(reader io.Reader, progress func()) {
	buffer := make([]byte, 32*1024)
	for {
		n, err := reader.Read(buffer)
		if n > 0 && progress != nil {
			progress()
		}
		if err != nil {
			return
		}
	}
}

func readTail(reader io.Reader, limit int) string {
	buffer := make([]byte, 32*1024)
	tail := make([]byte, 0, limit)
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			tail = append(tail, buffer[:n]...)
			if len(tail) > limit {
				tail = append([]byte(nil), tail[len(tail)-limit:]...)
			}
		}
		if err != nil {
			return string(tail)
		}
	}
}

func lastNonEmptyLine(value string) string {
	lines := strings.Split(value, "\n")
	for index := len(lines) - 1; index >= 0; index-- {
		if line := strings.TrimSpace(lines[index]); line != "" {
			return line
		}
	}
	return ""
}

func safeEnvironment(extra []string) []string {
	allowed := map[string]bool{
		"PATH": true, "HOME": true, "USER": true, "LOGNAME": true,
		"TMPDIR": true, "TMP": true, "TEMP": true,
		"LANG": true, "LC_ALL": true, "SHELL": true,
		"CODEX_HOME": true, "CLAUDE_CONFIG_DIR": true,
		"OPENAI_API_KEY": true, "ANTHROPIC_API_KEY": true,
	}
	environment := make([]string, 0, len(allowed)+len(extra)+1)
	for _, entry := range os.Environ() {
		name, _, ok := strings.Cut(entry, "=")
		if ok && allowed[name] {
			environment = append(environment, entry)
		}
	}
	environment = append(environment, "TERM=dumb")
	environment = append(environment, extra...)
	return environment
}

func Version(ctx context.Context, availability Availability) Availability {
	if !availability.Installed {
		return availability
	}
	command := exec.CommandContext(ctx, availability.Path, "--version")
	command.Env = safeEnvironment(nil)
	configureProcess(command)
	output, err := command.CombinedOutput()
	if err != nil {
		availability.Error = fmt.Sprintf("read version: %v", err)
		return availability
	}
	version := strings.TrimSpace(string(output))
	if len(version) > 160 {
		version = version[len(version)-160:]
	}
	availability.Version = version
	return availability
}

func IsCancellation(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
