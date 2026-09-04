package agent

import "fmt"

type claudeAdapter struct{}

func (claudeAdapter) Provider() Provider {
	return ProviderClaude
}

func (claudeAdapter) BinaryName() string {
	return "claude"
}

func (claudeAdapter) Build(request RunRequest) Invocation {
	allowedBash := fmt.Sprintf("Bash(%s cli *)", request.CLICommand)
	return Invocation{
		Path: request.BinaryPath,
		Args: []string{
			"-p",
			"--output-format", "stream-json",
			"--verbose",
			"--no-session-persistence",
			"--disable-slash-commands",
			"--no-chrome",
			"--setting-sources", "",
			"--tools", "Bash,Write",
			"--allowedTools", allowedBash + ",Write",
			"--permission-mode", "dontAsk",
		},
		Stdin: request.Prompt,
	}
}
