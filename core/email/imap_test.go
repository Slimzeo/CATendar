package email

import (
	"strings"
	"testing"
)

func TestReadMessageTextIgnoresAttachmentAndUsesPlainText(t *testing.T) {
	raw := strings.NewReader("MIME-Version: 1.0\r\n" +
		"Content-Type: multipart/mixed; boundary=test\r\n\r\n" +
		"--test\r\nContent-Type: text/plain; charset=utf-8\r\n\r\nMeeting tomorrow at 10.\r\n" +
		"--test\r\nContent-Type: application/pdf\r\nContent-Disposition: attachment; filename=x.pdf\r\n\r\nPDF DATA\r\n" +
		"--test--\r\n")

	text, err := readMessageText(raw)
	if err != nil {
		t.Fatalf("read message: %v", err)
	}
	if text != "Meeting tomorrow at 10." {
		t.Fatalf("unexpected text: %q", text)
	}
}

func TestSourceURLKeepsMessageIdentity(t *testing.T) {
	value := sourceURL("primary", 42, 7, "message@example.com")
	if !strings.HasPrefix(value, "catendar://email/primary?") || !strings.Contains(value, "uid=7") {
		t.Fatalf("unexpected source URL: %s", value)
	}
	ref, err := messageRefFromSource(value, "primary")
	if err != nil {
		t.Fatalf("parse source URL: %v", err)
	}
	if ref != "42:7" {
		t.Fatalf("unexpected message ref: %s", ref)
	}
}
