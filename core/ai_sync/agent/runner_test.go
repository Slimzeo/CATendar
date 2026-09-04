package agent

import (
	"strings"
	"testing"
	"time"
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

func TestPromptRequiresHighRecallAndReceiptAnchoredDeadlines(t *testing.T) {
	prompt := Prompt("./catendar-cli", time.Date(2026, time.September, 5, 12, 0, 0, 0, time.FixedZone("CST", 8*60*60)))

	for _, required := range []string{
		"MUST read",
		"interview",
		"面试",
		"receivedAt",
		"sentAt",
		"收到后三天内",
		"收到后 3 天内",
		"72 hours",
		"readError",
	} {
		if !strings.Contains(prompt, required) {
			t.Fatalf("prompt does not contain %q", required)
		}
	}
}
