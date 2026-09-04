package ai_sync

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	synccli "catendar-wails/core/ai_sync/cli"
)

const maxSessionRequestSize = 2 * 1024 * 1024

var sessionRequestIDPattern = regexp.MustCompile(`^[0-9a-f]{32}$`)

type localSession struct {
	directory   string
	token       string
	commandPath string
}

func newLocalSession(executablePath string) (*localSession, error) {
	directory, err := os.MkdirTemp("", "catendar-ai-sync-")
	if err != nil {
		return nil, fmt.Errorf("create AI Sync directory: %w", err)
	}
	if err := os.Chmod(directory, 0o700); err != nil {
		os.RemoveAll(directory)
		return nil, fmt.Errorf("secure AI Sync directory: %w", err)
	}
	for _, name := range []string{"requests", "responses"} {
		if err := os.Mkdir(filepath.Join(directory, name), 0o700); err != nil {
			os.RemoveAll(directory)
			return nil, fmt.Errorf("create AI Sync %s directory: %w", name, err)
		}
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		os.RemoveAll(directory)
		return nil, fmt.Errorf("create AI Sync capability: %w", err)
	}

	commandPath := shellQuote(executablePath)
	if runtime.GOOS != "windows" {
		linkPath := filepath.Join(directory, "catendar-cli")
		if err := os.Symlink(executablePath, linkPath); err == nil {
			commandPath = "./catendar-cli"
		}
	}
	return &localSession{
		directory:   directory,
		token:       hex.EncodeToString(tokenBytes),
		commandPath: commandPath,
	}, nil
}

func (s *localSession) Environment() []string {
	return []string{
		synccli.EnvSessionDir + "=" + s.directory,
		synccli.EnvSessionToken + "=" + s.token,
	}
}

func (s *localSession) WorkDir() string {
	return s.directory
}

func (s *localSession) CommandPath() string {
	return s.commandPath
}

func (s *localSession) Close() {
	os.RemoveAll(s.directory)
}

func (s *localSession) Serve(
	ctx context.Context,
	handler func(context.Context, synccli.Request) (any, error),
) error {
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			entries, err := os.ReadDir(filepath.Join(s.directory, "requests"))
			if err != nil {
				return fmt.Errorf("read AI Sync requests: %w", err)
			}
			for _, entry := range entries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
					continue
				}
				if err := s.handleRequest(ctx, entry.Name(), handler); err != nil {
					return err
				}
			}
		}
	}
}

func (s *localSession) handleRequest(
	ctx context.Context,
	name string,
	handler func(context.Context, synccli.Request) (any, error),
) error {
	requestPath := filepath.Join(s.directory, "requests", name)
	info, err := os.Lstat(requestPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect AI Sync request: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() > maxSessionRequestSize {
		os.Remove(requestPath)
		return nil
	}
	data, err := os.ReadFile(requestPath)
	if err != nil {
		return fmt.Errorf("read AI Sync request: %w", err)
	}
	os.Remove(requestPath)

	var request synccli.Request
	if err := json.Unmarshal(data, &request); err != nil ||
		!sessionRequestIDPattern.MatchString(request.ID) || name != request.ID+".json" {
		return nil
	}
	response := synccli.Response{OK: false}
	if !sameToken(request.Token, s.token) {
		response.Error = "AI Sync session is not authorized"
	} else {
		result, err := handler(ctx, request)
		if err != nil {
			response.Error = err.Error()
		} else {
			encoded, err := json.Marshal(result)
			if err != nil {
				response.Error = "CATendar Core could not encode the result"
			} else {
				response.OK = true
				response.Data = encoded
			}
		}
	}
	return s.writeResponse(request.ID, response)
}

func (s *localSession) writeResponse(id string, response synccli.Response) error {
	if !sessionRequestIDPattern.MatchString(id) {
		return errors.New("invalid AI Sync response id")
	}
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("encode AI Sync response: %w", err)
	}
	responseDirectory := filepath.Join(s.directory, "responses")
	path := filepath.Join(responseDirectory, id+".json")
	temporary, err := os.CreateTemp(responseDirectory, "."+id+"-*.tmp")
	if err != nil {
		return fmt.Errorf("create CLI response: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return fmt.Errorf("secure CLI response: %w", err)
	}
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return fmt.Errorf("write AI Sync response: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close AI Sync response: %w", err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("publish AI Sync response: %w", err)
	}
	return nil
}

func sameToken(left, right string) bool {
	if len(left) != len(right) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(left), []byte(right)) == 1
}

func shellQuote(value string) string {
	if runtime.GOOS == "windows" {
		return `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`
	}
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
