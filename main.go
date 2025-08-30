package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

type MailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

// getSecret returns a value either from environment variable or from a file path in *_FILE
func getSecret(envKey, fileKey string) string {
	if path := os.Getenv(fileKey); path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			return strings.TrimSpace(string(data))
		}
		log.Printf("Warning: failed to read secret file %s: %v", path, err)
	}
	return os.Getenv(envKey)
}

func sendMail(req MailRequest) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := getSecret("SMTP_PASS", "SMTP_PASS_FILE")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s", req.From, strings.Join(req.To, ","), req.Subject, req.Body))

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, req.From, req.To, msg)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Only POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// Token auth
	expectedToken := getSecret("API_TOKEN", "API_TOKEN_FILE")
	authHeader := r.Header.Get("Authorization")
	if expectedToken != "" {
		if !strings.HasPrefix(authHeader, "Bearer ") || strings.TrimPrefix(authHeader, "Bearer ") != expectedToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Decode JSON
	var req MailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Send mail
	if err := sendMail(req); err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent"))
}

func main() {
	http.HandleFunc("/v1/send", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
