package ai_sync

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	synccli "catendar-wails/core/ai_sync/cli"
)

func TestLocalSessionCarriesCLIRequestToCore(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("resolve executable: %v", err)
	}
	session, err := newLocalSession(executable)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	t.Cleanup(session.Close)
	t.Setenv(synccli.EnvSessionDir, session.directory)
	t.Setenv(synccli.EnvSessionToken, session.token)

	ctx, cancel := context.WithCancel(context.Background())
	serverDone := make(chan error, 1)
	go func() {
		serverDone <- session.Serve(ctx, func(_ context.Context, request synccli.Request) (any, error) {
			if request.Command != "email.list" {
				t.Fatalf("unexpected command: %s", request.Command)
			}
			return map[string]int{"total": 3}, nil
		})
	}()

	var stdout, stderr bytes.Buffer
	if code := synccli.Run(context.Background(), []string{"email", "list"}, nil, &stdout, &stderr); code != 0 {
		t.Fatalf("CLI failed with %d: %s", code, stderr.String())
	}
	var response synccli.Response
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !response.OK || string(response.Data) != `{"total":3}` {
		t.Fatalf("unexpected response: %s", stdout.String())
	}

	cancel()
	if err := <-serverDone; !errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected server result: %v", err)
	}
}

func TestLocalSessionRejectsUnsafeResponseID(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("resolve executable: %v", err)
	}
	session, err := newLocalSession(executable)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	t.Cleanup(session.Close)

	if err := session.writeResponse("../../outside", synccli.Response{OK: true}); err == nil {
		t.Fatal("unsafe response id was accepted")
	}
}
