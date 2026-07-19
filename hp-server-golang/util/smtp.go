package util

import (
	"crypto/tls"
	"fmt"
	"hp-server-lib/config"
	"net/smtp"
)

func SendEmail(to, subject, body string) error {
	cfg := config.ConfigData.System.Smtp
	if !cfg.Enabled {
		return nil
	}

	from := cfg.From
	if cfg.FromName != "" {
		from = fmt.Sprintf("%s <%s>", cfg.FromName, cfg.From)
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		from, to, subject, body)

	var auth smtp.Auth
	if cfg.Username != "" && cfg.Password != "" {
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	if cfg.EnableSSL {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         cfg.Host,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return err
		}

		client, err := smtp.NewClient(conn, cfg.Host)
		if err != nil {
			return err
		}

		if auth != nil {
			if err := client.Auth(auth); err != nil {
				return err
			}
		}

		if err := client.Mail(cfg.From); err != nil {
			return err
		}

		if err := client.Rcpt(to); err != nil {
			return err
		}

		w, err := client.Data()
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(msg))
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}

		return client.Quit()
	} else {
		return smtp.SendMail(addr, auth, cfg.From, []string{to}, []byte(msg))
	}
}

func SendRegisterEmail(email, username string) error {
	subject := "欢迎注册HP-Lite内网穿透"
	body := fmt.Sprintf(`<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
	<h2 style="color: #4b6ff6;">欢迎注册HP-Lite内网穿透</h2>
	<p>尊敬的 %s，</p>
	<p>感谢您注册HP-Lite内网穿透服务！</p>
	<p>您的用户名：%s</p>
	<p>如有任何问题，请联系管理员。</p>
	<p>祝您使用愉快！</p>
	<hr>
	<p style="color: #888; font-size: 12px;">HP-Lite内网穿透团队</p>
</div>`, username, username)
	return SendEmail(email, subject, body)
}
