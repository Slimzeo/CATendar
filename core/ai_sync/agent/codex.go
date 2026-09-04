package agent

type codexAdapter struct{}

func (codexAdapter) Provider() Provider {
	return ProviderCodex
}

func (codexAdapter) BinaryName() string {
	return "codex"
}

func (codexAdapter) Build(request RunRequest) Invocation {
	return Invocation{
		Path: request.BinaryPath,
		Args: []string{
			"exec",
			"--ephemeral",
			"--ignore-user-config",
			"--ignore-rules",
			"--skip-git-repo-check",
			"--sandbox", "workspace-write",
			"--json",
			"--color", "never",
			"-c", `approval_policy="never"`,
			"-C", request.WorkDir,
			"-",
		},
		Stdin: request.Prompt,
	}
}
