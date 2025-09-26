package utils

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

// Verifica si una string contiene únicamente dígitos
// Retorna un boolean
func isAllDigits(s string) bool {
	for _, c := range s {
		if c < 48 || c > 57 {
			return false
		}
	}
	return true
}

// Retorna un mensaje de error legible para cada validador
// Recibe el nombre como parámetro
func getTagMessage(tag string) string {
	switch tag {
	case "phone":
		return "Phone number must contain digits only"

	case "password":
		return "Password must be at least 8 characters long"

	case "email":
		return "Invalid email address"

	case "plate":
		return "Invalid license plate format"
	}

	return "Invalid request format"
}

// Envía un correo electrónico vía Gmail con contenido HTML
// Recibe el recipiente y el contenido del correo como parámetros
// Puede retornar un error
func SendEmail(receiver, subject, htmlContent string) error {
	sender := "deaquiparalla.gt@gmail.com"
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	conn, err := net.Dial("tcp", net.JoinHostPort(smtpHost, fmt.Sprintf("%d", smtpPort)))
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{ServerName: smtpHost}
	if err = client.StartTLS(tlsConfig); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", sender, EnvMustGet("EMAIL_PASSWORD"), smtpHost)
	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(sender); err != nil {
		return err
	}
	if err = client.Rcpt(receiver); err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["From"] = sender
	headers["To"] = receiver
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlContent

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	w.Close()
	client.Quit()
	return nil
}
