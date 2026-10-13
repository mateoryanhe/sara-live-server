package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	sysentity "xr-game-server/entity/sys"
)

func cfEmailSMTPSettings(cfg *sysentity.CfEmailCfg) (host string, port int, username, password, from string, ok bool) {
	if cfg == nil {
		return "", 0, "", "", "", false
	}
	host = strings.TrimSpace(cfg.SmtpHost)
	username = strings.TrimSpace(cfg.SmtpUsername)
	password = strings.TrimSpace(cfg.SmtpPassword)
	from = strings.TrimSpace(cfg.FromEmail)
	port = cfg.SmtpPort
	if port <= 0 {
		port = 587
	}
	ok = host != "" && username != "" && password != "" && from != ""
	return host, port, username, password, from, ok
}

func sendEmailViaSMTP(ctx context.Context, cfg *sysentity.CfEmailCfg, toEmail, subject, plainBody, htmlBody string) error {
	host, port, username, password, from, ok := cfEmailSMTPSettings(cfg)
	if !ok {
		return fmt.Errorf("smtp cfg incomplete")
	}
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	auth := smtp.PlainAuth("", username, password, host)
	msg := formatVerificationEmail(from, toEmail, subject, plainBody, htmlBody)
	return sendMailStartTLS(ctx, addr, host, auth, from, []string{toEmail}, msg)
}

func formatVerificationEmail(from, to, subject, plainBody, htmlBody string) []byte {
	var buf bytes.Buffer
	buf.WriteString("From: " + from + "\r\n")
	buf.WriteString("To: " + to + "\r\n")
	buf.WriteString("Subject: " + mime.QEncoding.Encode("utf-8", subject) + "\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")
	if strings.TrimSpace(htmlBody) == "" {
		buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		buf.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
		buf.WriteString(plainBody)
		return buf.Bytes()
	}
	boundary := randomMIMEBoundary()
	buf.WriteString("Content-Type: multipart/alternative; boundary=" + boundary + "\r\n\r\n")
	writeMIMEPart(&buf, boundary, "text/plain; charset=UTF-8", plainBody)
	writeMIMEPart(&buf, boundary, "text/html; charset=UTF-8", htmlBody)
	buf.WriteString("--" + boundary + "--\r\n")
	return buf.Bytes()
}

func randomMIMEBoundary() string {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "xr-email-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return "xr-email-" + hex.EncodeToString(b[:])
}

func writeMIMEPart(buf *bytes.Buffer, boundary, contentType, body string) {
	buf.WriteString("--" + boundary + "\r\n")
	buf.WriteString("Content-Type: " + contentType + "\r\n")
	buf.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
	buf.WriteString(body)
	buf.WriteString("\r\n")
}

func sendMailStartTLS(ctx context.Context, addr, serverName string, auth smtp.Auth, from string, to []string, msg []byte) error {
	dialer := net.Dialer{Timeout: 30 * time.Second}
	if deadline, ok := ctx.Deadline(); ok {
		dialer.Deadline = deadline
	}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, serverName)
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(&tls.Config{ServerName: serverName}); err != nil {
			return err
		}
	}
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}
	if err = client.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range to {
		if err = client.Rcpt(rcpt); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write(msg); err != nil {
		_ = w.Close()
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	return client.Quit()
}
