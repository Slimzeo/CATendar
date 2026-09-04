package email

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

type Service struct {
	db          *sql.DB
	credentials CredentialStore
	reader      Reader
}

func NewService(db *sql.DB, credentials CredentialStore, reader Reader) *Service {
	return &Service{db: db, credentials: credentials, reader: reader}
}

func (s *Service) GetAccount(ctx context.Context) (Account, error) {
	account := Account{
		ID:       DefaultAccountID,
		IMAPPort: 993,
		Folder:   "INBOX",
		UseTLS:   true,
	}
	var useTLS int
	err := s.db.QueryRowContext(ctx, `
		SELECT address, username, imap_host, imap_port, folder, use_tls
		FROM email_accounts WHERE id = ?
	`, DefaultAccountID).Scan(
		&account.Address,
		&account.Username,
		&account.IMAPHost,
		&account.IMAPPort,
		&account.Folder,
		&useTLS,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return account, nil
	}
	if err != nil {
		return Account{}, fmt.Errorf("read email account: %w", err)
	}

	account.UseTLS = useTLS == 1
	account.Configured = true
	hasCredential, err := s.credentials.Exists(account.ID)
	if err != nil {
		return Account{}, err
	}
	account.HasCredential = hasCredential
	account.Configured = hasCredential
	return account, nil
}

func (s *Service) SaveAccount(ctx context.Context, input AccountInput) (Account, error) {
	account, err := normalizeAccount(input)
	if err != nil {
		return Account{}, err
	}

	if input.Secret != "" {
		if err := s.credentials.Set(account.ID, input.Secret); err != nil {
			return Account{}, err
		}
	} else if exists, err := s.credentials.Exists(account.ID); err != nil {
		return Account{}, err
	} else if !exists {
		return Account{}, errors.New("email password or app authorization code is required")
	}

	if _, err := s.db.ExecContext(ctx, `
		INSERT INTO email_accounts (id, address, username, imap_host, imap_port, folder, use_tls, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			address = excluded.address,
			username = excluded.username,
			imap_host = excluded.imap_host,
			imap_port = excluded.imap_port,
			folder = excluded.folder,
			use_tls = excluded.use_tls,
			updated_at = excluded.updated_at
	`, account.ID, account.Address, account.Username, account.IMAPHost, account.IMAPPort,
		account.Folder, boolToInt(account.UseTLS), time.Now()); err != nil {
		return Account{}, fmt.Errorf("save email account: %w", err)
	}
	account.Configured = true
	account.HasCredential = true
	return account, nil
}

func (s *Service) TestAccount(ctx context.Context, input AccountInput) error {
	account, err := normalizeAccount(input)
	if err != nil {
		return err
	}
	secret := input.Secret
	if secret == "" {
		secret, err = s.credentials.Get(account.ID)
		if err != nil {
			return err
		}
	}
	return s.reader.Test(ctx, account, secret)
}

func (s *Service) ListRecent(ctx context.Context, since time.Time, limit int) (Account, []MessageSummary, error) {
	account, secret, err := s.configuredAccount(ctx)
	if err != nil {
		return Account{}, nil, err
	}
	messages, err := s.reader.ListRecent(ctx, account, secret, since, limit)
	if err != nil {
		return Account{}, nil, err
	}
	return account, messages, nil
}

func (s *Service) ReadMessages(ctx context.Context, ids []string) ([]Message, error) {
	account, secret, err := s.configuredAccount(ctx)
	if err != nil {
		return nil, err
	}
	return s.reader.Read(ctx, account, secret, ids)
}

// ReadMessagesForAccount keeps an AI Sync run pinned to the mailbox metadata
// it started with, even if the UI changes its saved account concurrently.
func (s *Service) ReadMessagesForAccount(ctx context.Context, account Account, ids []string) ([]Message, error) {
	if !account.Configured || account.ID != DefaultAccountID {
		return nil, errors.New("email account is not configured")
	}
	secret, err := s.credentials.Get(account.ID)
	if err != nil {
		return nil, err
	}
	return s.reader.Read(ctx, account, secret, ids)
}

func (s *Service) ReadSource(ctx context.Context, source string) (*Message, error) {
	account, secret, err := s.configuredAccount(ctx)
	if err != nil {
		return nil, err
	}
	id, err := messageRefFromSource(source, account.ID)
	if err != nil {
		return nil, err
	}
	messages, err := s.reader.Read(ctx, account, secret, []string{id})
	if err != nil {
		return nil, err
	}
	if len(messages) != 1 {
		return nil, errors.New("source email was not found")
	}
	if messages[0].ReadError != "" {
		return nil, errors.New(messages[0].ReadError)
	}
	return &messages[0], nil
}

func (s *Service) configuredAccount(ctx context.Context) (Account, string, error) {
	account, err := s.GetAccount(ctx)
	if err != nil {
		return Account{}, "", err
	}
	if !account.Configured {
		return Account{}, "", errors.New("email account is not configured")
	}
	secret, err := s.credentials.Get(account.ID)
	if err != nil {
		return Account{}, "", err
	}
	return account, secret, nil
}

func normalizeAccount(input AccountInput) (Account, error) {
	parsedAddress, err := mail.ParseAddress(strings.TrimSpace(input.Address))
	if err != nil || parsedAddress.Address == "" {
		return Account{}, errors.New("enter a valid email address")
	}
	host := strings.TrimSpace(input.IMAPHost)
	if host == "" {
		return Account{}, errors.New("IMAP host is required")
	}
	if input.IMAPPort < 1 || input.IMAPPort > 65535 {
		return Account{}, errors.New("IMAP port must be between 1 and 65535")
	}
	username := strings.TrimSpace(input.Username)
	if username == "" {
		username = parsedAddress.Address
	}
	folder := strings.TrimSpace(input.Folder)
	if folder == "" {
		folder = "INBOX"
	}
	return Account{
		ID:       DefaultAccountID,
		Address:  parsedAddress.Address,
		Username: username,
		IMAPHost: host,
		IMAPPort: input.IMAPPort,
		Folder:   folder,
		UseTLS:   input.UseTLS,
	}, nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
