package mailer

import (
	"fmt"
	"net/smtp"

	"goose-neon-demo/config"
)

type Mailer struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Mailer {
	return &Mailer{cfg: cfg}
}

func (m *Mailer) SendMagicLink(toEmail, verifyURL string) error {
	subject := "Your sign-in link"
	body := fmt.Sprintf(
		"Click the link below to sign in. This link expires shortly and can only be used once.\r\n\r\n%s\r\n\r\nIf you didn't request this, you can ignore this email.",
		verifyURL,
	)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.cfg.SMTPFrom, toEmail, subject, body,
	)

	addr := fmt.Sprintf("%s:%s", m.cfg.SMTPHost, m.cfg.SMTPPort)
	auth := smtp.PlainAuth("", m.cfg.SMTPUsername, m.cfg.SMTPPassword, m.cfg.SMTPHost)

	return smtp.SendMail(addr, auth, m.cfg.SMTPFrom, []string{toEmail}, []byte(msg))
}
