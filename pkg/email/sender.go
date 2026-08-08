package email

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Pass     string
	From     string
	FromName string
}

type Sender struct {
	cfg Config
}

func NewSender(cfg Config) *Sender {
	return &Sender{cfg: cfg}
}

func (s *Sender) Configured() bool {
	return s.cfg.Host != "" && s.cfg.Port != ""
}

func (s *Sender) FromAddress() string {
	if s.cfg.From != "" {
		return s.cfg.From
	}
	return s.cfg.User
}

func (s *Sender) FromHeader() string {
	addr := s.FromAddress()
	name := s.cfg.FromName
	if name == "" {
		return addr
	}
	return fmt.Sprintf("%s <%s>", mime.QEncoding.Encode("utf-8", name), addr)
}

// Send delivers a plain-text email to a single recipient.
func (s *Sender) Send(to, subject, body string) error {
	if !s.Configured() {
		return fmt.Errorf("smtp not configured")
	}

	addr := net.JoinHostPort(s.cfg.Host, s.cfg.Port)
	msg := buildMessage(s.FromHeader(), to, subject, body)
	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Pass, s.cfg.Host)

	from := s.cfg.User
	if s.cfg.From != "" {
		from = s.cfg.From
	}

	if s.cfg.Port == "465" {
		return sendTLS(addr, from, []string{to}, []byte(msg), auth)
	}

	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}

func buildMessage(from, to, subject, body string) string {
	var b strings.Builder
	b.WriteString("From: ")
	b.WriteString(from)
	b.WriteString("\r\n")
	b.WriteString("To: ")
	b.WriteString(to)
	b.WriteString("\r\n")
	b.WriteString("Subject: ")
	b.WriteString(mime.QEncoding.Encode("utf-8", subject))
	b.WriteString("\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	b.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	b.WriteString("Date: ")
	b.WriteString(time.Now().Format(time.RFC1123Z))
	b.WriteString("\r\n")
	b.WriteString("Message-ID: <")
	b.WriteString(uuid.NewString())
	b.WriteString("@")
	b.WriteString("foundry")
	b.WriteString(">\r\n")
	b.WriteString("\r\n")
	b.WriteString(body)
	return b.String()
}

func sendTLS(addr, from string, to []string, msg []byte, auth smtp.Auth) error {
	host, _, _ := net.SplitHostPort(addr)
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: host,
		MinVersion: tls.VersionTLS12,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return err
	}
	if err := client.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range to {
		if err := client.Rcpt(rcpt); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return client.Quit()
}
