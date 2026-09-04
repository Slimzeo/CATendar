package agent

import (
	"strings"
	"testing"
)

func TestReadTailDrainsLargeOutputAndKeepsTheEnd(t *testing.T) {
	input := strings.Repeat("a", 2*1024*1024) + "\nfinal error\n"
	tail := readTail(strings.NewReader(input), 1024)
	if len(tail) > 1024 {
		t.Fatalf("tail exceeds limit: %d", len(tail))
	}
	if got := lastNonEmptyLine(tail); got != "final error" {
		t.Fatalf("unexpected final line: %q", got)
	}
}

func TestAdaptersUseTheSharedCATendarCLI(t *testing.T) {
	request := RunRequest{
		BinaryPath: "/usr/local/bin/agent",
		WorkDir:    "/tmp/catendar-run",
		CLICommand: "./catendar-cli",
		Prompt:     "plan",
	}
	for _, provider := range []Provider{ProviderCodex, ProviderClaude} {
		adapter, err := ForProvider(provider)
		if err != nil {
			t.Fatalf("resolve %s adapter: %v", provider, err)
		}
		invocation := adapter.Build(request)
		if invocation.Path != request.BinaryPath || invocation.Stdin != request.Prompt {
			t.Fatalf("unexpected %s invocation: %#v", provider, invocation)
		}
		if provider == ProviderClaude && !strings.Contains(strings.Join(invocation.Args, " "), request.CLICommand+" cli") {
			t.Fatalf("Claude invocation does not constrain Bash to the shared CLI: %#v", invocation.Args)
		}
	}
}
