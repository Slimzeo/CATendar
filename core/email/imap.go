package email

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"html"
	"io"
	"net"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
	messageMail "github.com/emersion/go-message/mail"
)

const (
	imapTimeout        = 20 * time.Second
	maxMessageBodySize = 256 * 1024
)

var (
	unsafeHTMLPattern = regexp.MustCompile(`(?is)<(script|style)[^>]*>.*?</(script|style)>`)
	htmlTagPattern    = regexp.MustCompile(`(?s)<[^>]+>`)
	manySpacesPattern = regexp.MustCompile(`[\t ]+`)
	manyLinesPattern  = regexp.MustCompile(`\n{3,}`)
)

type IMAPReader struct{}

func NewIMAPReader() IMAPReader {
	return IMAPReader{}
}

func (IMAPReader) Test(ctx context.Context, account Account, secret string) error {
	connection, err := connect(ctx, account, secret)
	if err != nil {
		return err
	}
	defer connection.Logout()
	if _, err := connection.Select(account.Folder, true); err != nil {
		return fmt.Errorf("open IMAP folder %q: %w", account.Folder, err)
	}
	return nil
}

func (IMAPReader) ListRecent(
	ctx context.Context,
	account Account,
	secret string,
	since time.Time,
	limit int,
) ([]MessageSummary, error) {
	connection, err := connect(ctx, account, secret)
	if err != nil {
		return nil, err
	}
	defer connection.Logout()

	mailbox, err := connection.Select(account.Folder, true)
	if err != nil {
		return nil, fmt.Errorf("open IMAP folder %q: %w", account.Folder, err)
	}
	criteria := imap.NewSearchCriteria()
	criteria.Since = since
	uids, err := connection.UidSearch(criteria)
	if err != nil {
		return nil, fmt.Errorf("search recent email: %w", err)
	}
	if len(uids) == 0 {
		return []MessageSummary{}, nil
	}
	if limit <= 0 {
		limit = 200
	}
	if len(uids) > limit {
		uids = uids[len(uids)-limit:]
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)
	messages := make(chan *imap.Message, len(uids))
	fetchDone := make(chan error, 1)
	go func() {
		fetchDone <- connection.UidFetch(seqSet, []imap.FetchItem{
			imap.FetchUid,
			imap.FetchEnvelope,
			imap.FetchInternalDate,
		}, messages)
	}()

	result := make([]MessageSummary, 0, len(uids))
	for item := range messages {
		if item == nil || item.Envelope == nil {
			continue
		}
		receivedAt := item.InternalDate
		if !item.Envelope.Date.IsZero() {
			receivedAt = item.Envelope.Date
		}
		messageID := canonicalMessageID(item.Envelope.MessageId, mailbox.UidValidity, item.Uid)
		result = append(result, MessageSummary{
			ID:         formatMessageRef(mailbox.UidValidity, item.Uid),
			MessageID:  messageID,
			Subject:    strings.TrimSpace(item.Envelope.Subject),
			Sender:     formatAddresses(item.Envelope.From),
			ReceivedAt: receivedAt,
			SourceURL:  sourceURL(account.ID, mailbox.UidValidity, item.Uid, messageID),
		})
	}
	if err := <-fetchDone; err != nil {
		return nil, fmt.Errorf("read email headers: %w", err)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ReceivedAt.After(result[j].ReceivedAt)
	})
	return result, nil
}

func (IMAPReader) Read(
	ctx context.Context,
	account Account,
	secret string,
	ids []string,
) ([]Message, error) {
	if len(ids) == 0 {
		return []Message{}, nil
	}
	if len(ids) > 20 {
		return nil, errors.New("read at most 20 email messages per request")
	}

	connection, err := connect(ctx, account, secret)
	if err != nil {
		return nil, err
	}
	defer connection.Logout()
	mailbox, err := connection.Select(account.Folder, true)
	if err != nil {
		return nil, fmt.Errorf("open IMAP folder %q: %w", account.Folder, err)
	}

	seqSet := new(imap.SeqSet)
	requested := make(map[uint32]string, len(ids))
	for _, id := range ids {
		uidValidity, uid, err := parseMessageRef(id)
		if err != nil {
			return nil, err
		}
		if uidValidity != mailbox.UidValidity {
			return nil, errors.New("email mailbox identity changed; run AI Sync again")
		}
		seqSet.AddNum(uid)
		requested[uid] = id
	}

	section := &imap.BodySectionName{Peek: true}
	messages := make(chan *imap.Message, len(ids))
	fetchDone := make(chan error, 1)
	go func() {
		fetchDone <- connection.UidFetch(seqSet, []imap.FetchItem{
			imap.FetchUid,
			imap.FetchEnvelope,
			imap.FetchInternalDate,
			section.FetchItem(),
		}, messages)
	}()

	result := make([]Message, 0, len(ids))
	for item := range messages {
		id, ok := requested[item.Uid]
		if !ok || item.Envelope == nil {
			continue
		}
		body := item.GetBody(section)
		if body == nil {
			continue
		}
		text, err := readMessageText(body)
		if err != nil {
			return nil, fmt.Errorf("parse email %s: %w", id, err)
		}
		receivedAt := item.InternalDate
		if !item.Envelope.Date.IsZero() {
			receivedAt = item.Envelope.Date
		}
		messageID := canonicalMessageID(item.Envelope.MessageId, mailbox.UidValidity, item.Uid)
		result = append(result, Message{
			MessageSummary: MessageSummary{
				ID:         id,
				MessageID:  messageID,
				Subject:    strings.TrimSpace(item.Envelope.Subject),
				Sender:     formatAddresses(item.Envelope.From),
				ReceivedAt: receivedAt,
				SourceURL:  sourceURL(account.ID, mailbox.UidValidity, item.Uid, messageID),
			},
			Text: text,
		})
	}
	if err := <-fetchDone; err != nil {
		return nil, fmt.Errorf("read email messages: %w", err)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ReceivedAt.After(result[j].ReceivedAt)
	})
	return result, nil
}

func connect(ctx context.Context, account Account, secret string) (*client.Client, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	address := net.JoinHostPort(account.IMAPHost, strconv.Itoa(account.IMAPPort))
	dialer := imapDialer{ctx: ctx}
	var connection *client.Client
	var err error
	if account.UseTLS {
		connection, err = client.DialWithDialerTLS(dialer, address, &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: account.IMAPHost,
		})
	} else {
		connection, err = client.DialWithDialer(dialer, address)
	}
	if err != nil {
		return nil, fmt.Errorf("connect to IMAP server: %w", err)
	}
	connection.Timeout = imapTimeout
	if err := connection.Login(account.Username, secret); err != nil {
		connection.Logout()
		return nil, fmt.Errorf("authenticate with IMAP server: %w", err)
	}
	return connection, nil
}

type imapDialer struct {
	ctx context.Context
}

func (d imapDialer) Dial(network, address string) (net.Conn, error) {
	dialer := net.Dialer{Timeout: imapTimeout}
	connection, err := dialer.DialContext(d.ctx, network, address)
	if err != nil {
		return nil, err
	}
	if err := connection.SetDeadline(time.Now().Add(imapTimeout)); err != nil {
		connection.Close()
		return nil, err
	}
	return connection, nil
}

func readMessageText(reader io.Reader) (string, error) {
	mailReader, err := messageMail.CreateReader(io.LimitReader(reader, maxMessageBodySize))
	if err != nil && !message.IsUnknownCharset(err) {
		return "", err
	}

	var plainParts []string
	var htmlParts []string
	for {
		part, err := mailReader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil && !message.IsUnknownCharset(err) {
			return "", err
		}
		header, ok := part.Header.(*messageMail.InlineHeader)
		if !ok {
			continue
		}
		contentType, _, _ := header.ContentType()
		content, err := io.ReadAll(io.LimitReader(part.Body, maxMessageBodySize))
		if err != nil {
			return "", err
		}
		switch strings.ToLower(contentType) {
		case "text/plain", "":
			plainParts = append(plainParts, string(content))
		case "text/html":
			htmlParts = append(htmlParts, htmlToText(string(content)))
		}
	}

	text := strings.Join(plainParts, "\n\n")
	if strings.TrimSpace(text) == "" {
		text = strings.Join(htmlParts, "\n\n")
	}
	return cleanText(text), nil
}

func htmlToText(value string) string {
	value = unsafeHTMLPattern.ReplaceAllString(value, " ")
	value = strings.NewReplacer(
		"<br>", "\n", "<br/>", "\n", "<br />", "\n",
		"</p>", "\n", "</div>", "\n", "</li>", "\n",
	).Replace(value)
	return html.UnescapeString(htmlTagPattern.ReplaceAllString(value, " "))
}

func cleanText(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	lines := strings.Split(value, "\n")
	for index := range lines {
		lines[index] = strings.TrimSpace(manySpacesPattern.ReplaceAllString(lines[index], " "))
	}
	value = strings.TrimSpace(strings.Join(lines, "\n"))
	value = manyLinesPattern.ReplaceAllString(value, "\n\n")
	if len(value) > maxMessageBodySize {
		value = value[:maxMessageBodySize]
	}
	return value
}

func formatAddresses(addresses []*imap.Address) string {
	values := make([]string, 0, len(addresses))
	for _, address := range addresses {
		if address == nil {
			continue
		}
		mailbox := address.Address()
		if address.PersonalName != "" {
			values = append(values, fmt.Sprintf("%s <%s>", address.PersonalName, mailbox))
		} else {
			values = append(values, mailbox)
		}
	}
	return strings.Join(values, ", ")
}

func canonicalMessageID(messageID string, uidValidity, uid uint32) string {
	messageID = strings.Trim(strings.TrimSpace(messageID), "<>")
	if messageID != "" {
		return messageID
	}
	return fmt.Sprintf("imap-%d-%d", uidValidity, uid)
}

func formatMessageRef(uidValidity, uid uint32) string {
	return fmt.Sprintf("%d:%d", uidValidity, uid)
}

func parseMessageRef(value string) (uint32, uint32, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid email id")
	}
	uidValidity, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return 0, 0, errors.New("invalid email id")
	}
	uid, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return 0, 0, errors.New("invalid email id")
	}
	return uint32(uidValidity), uint32(uid), nil
}

func sourceURL(accountID string, uidValidity, uid uint32, messageID string) string {
	query := url.Values{}
	query.Set("uidValidity", strconv.FormatUint(uint64(uidValidity), 10))
	query.Set("uid", strconv.FormatUint(uint64(uid), 10))
	query.Set("messageId", messageID)
	return (&url.URL{
		Scheme:   "catendar",
		Host:     "email",
		Path:     "/" + accountID,
		RawQuery: query.Encode(),
	}).String()
}

func messageRefFromSource(value, accountID string) (string, error) {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme != "catendar" || parsed.Host != "email" {
		return "", errors.New("invalid CATendar email link")
	}
	if strings.TrimPrefix(parsed.Path, "/") != accountID {
		return "", errors.New("email link belongs to another account")
	}
	uidValidity, err := strconv.ParseUint(parsed.Query().Get("uidValidity"), 10, 32)
	if err != nil {
		return "", errors.New("email link has an invalid mailbox identity")
	}
	uid, err := strconv.ParseUint(parsed.Query().Get("uid"), 10, 32)
	if err != nil {
		return "", errors.New("email link has an invalid message identity")
	}
	return formatMessageRef(uint32(uidValidity), uint32(uid)), nil
}
