package email

import (
	"context"
	"time"
)

const DefaultAccountID = "primary"

type Account struct {
	ID            string `json:"id"`
	Address       string `json:"address"`
	Username      string `json:"username"`
	IMAPHost      string `json:"imapHost"`
	IMAPPort      int    `json:"imapPort"`
	Folder        string `json:"folder"`
	UseTLS        bool   `json:"useTLS"`
	Configured    bool   `json:"configured"`
	HasCredential bool   `json:"hasCredential"`
}

type AccountInput struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	IMAPHost string `json:"imapHost"`
	IMAPPort int    `json:"imapPort"`
	Folder   string `json:"folder"`
	UseTLS   bool   `json:"useTLS"`
	Secret   string `json:"secret"`
}

type MessageSummary struct {
	ID         string     `json:"id"`
	MessageID  string     `json:"messageId"`
	Subject    string     `json:"subject"`
	Sender     string     `json:"sender"`
	ReceivedAt time.Time  `json:"receivedAt"`
	SentAt     *time.Time `json:"sentAt,omitempty"`
	SourceURL  string     `json:"sourceUrl"`
}

type Message struct {
	MessageSummary
	Text      string `json:"text"`
	ReadError string `json:"readError,omitempty"`
}

type Reader interface {
	Test(ctx context.Context, account Account, secret string) error
	ListRecent(ctx context.Context, account Account, secret string, since time.Time, limit int) ([]MessageSummary, error)
	Read(ctx context.Context, account Account, secret string, ids []string) ([]Message, error)
}
