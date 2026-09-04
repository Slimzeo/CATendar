package email

import (
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/emersion/go-imap"
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

func TestReadMessageTextDoesNotTruncateLargeMultipartMessage(t *testing.T) {
	raw := strings.NewReader("MIME-Version: 1.0\r\n" +
		"Content-Type: multipart/mixed; boundary=test\r\n\r\n" +
		"--test\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n" +
		"ByteDance interview is Tuesday at 10:00.\r\n" +
		"--test\r\nContent-Type: application/octet-stream\r\n" +
		"Content-Disposition: attachment; filename=large.bin\r\n\r\n" +
		strings.Repeat("A", maxMessageTextSize+32*1024) + "\r\n" +
		"--test--\r\n")

	text, err := readMessageText(raw)
	if err != nil {
		t.Fatalf("read large message: %v", err)
	}
	if text != "ByteDance interview is Tuesday at 10:00." {
		t.Fatalf("unexpected text: %q", text)
	}
}

func TestMessageSummaryKeepsSentAndReceivedTimesDistinct(t *testing.T) {
	receivedAt := time.Date(2026, time.September, 3, 16, 44, 32, 0, time.FixedZone("CST", 8*60*60))
	sentAt := receivedAt.Add(-9 * time.Minute)
	item := &imap.Message{
		Uid:          2066,
		InternalDate: receivedAt,
		Envelope: &imap.Envelope{
			Date:      sentAt,
			MessageId: "deadline@example.com",
			Subject:   "Complete within three days of receipt",
		},
	}

	summary := messageSummary(item, "primary", 1763350974)
	if !summary.ReceivedAt.Equal(receivedAt) {
		t.Fatalf("receivedAt = %s, want %s", summary.ReceivedAt, receivedAt)
	}
	if !summary.SentAt.Equal(sentAt) {
		t.Fatalf("sentAt = %s, want %s", summary.SentAt, sentAt)
	}
}

func TestCleanTextKeepsUTF8ValidAtLimit(t *testing.T) {
	text := cleanText(strings.Repeat("日", maxMessageTextSize))
	if len(text) > maxMessageTextSize {
		t.Fatalf("cleaned text exceeds limit: %d", len(text))
	}
	if !utf8.ValidString(text) {
		t.Fatal("cleaned text ends with an incomplete UTF-8 sequence")
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
