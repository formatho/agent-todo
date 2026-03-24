package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"time"

	"github.com/formatho/agent-todo/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EmailConfig holds email service configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	UseTLS       bool
}

// EmailService handles email sending
type EmailService struct {
	config EmailConfig
	db     *gorm.DB
}

// NewEmailService creates a new email service
func NewEmailService(db *gorm.DB) *EmailService {
	config := EmailConfig{
		SMTPHost:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvOrDefault("SMTP_PORT", "587"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    getEnvOrDefault("FROM_EMAIL", "noreply@formatho.com"),
		FromName:     getEnvOrDefault("FROM_NAME", "Formatho Agent Todo"),
		UseTLS:       os.Getenv("SMTP_USE_TLS") != "false",
	}

	return &EmailService{
		config: config,
		db:     db,
	}
}

// SendEmail sends an email immediately
func (s *EmailService) SendEmail(to, subject, bodyHTML, bodyText string) error {
	if s.config.SMTPUsername == "" || s.config.SMTPPassword == "" {
		return fmt.Errorf("SMTP credentials not configured")
	}

	// Build the email message
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	msg := s.buildMessage(from, to, subject, bodyHTML, bodyText)

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	
	var auth smtp.Auth
	if s.config.SMTPUsername != "" {
		auth = smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
	}

	// Send email
	if s.config.UseTLS {
		return s.sendWithTLS(addr, auth, from, []string{to}, []byte(msg))
	}

	return smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, []byte(msg))
}

// sendWithTLS sends email with TLS
func (s *EmailService) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.SMTPHost,
	})
	if err != nil {
		// Fallback to STARTTLS
		return smtp.SendMail(addr, auth, from, to, msg)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.SMTPHost)
	if err != nil {
		return err
	}
	defer client.Close()

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	return w.Close()
}

// buildMessage constructs the email message
func (s *EmailService) buildMessage(from, to, subject, bodyHTML, bodyText string) string {
	msg := fmt.Sprintf("From: %s\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-version: 1.0;\r\n"
	msg += "Content-Type: multipart/alternative; boundary=\"boundary\"\r\n"
	msg += "\r\n"
	msg += "--boundary\r\n"
	msg += "Content-Type: text/plain; charset=\"UTF-8\"\r\n"
	msg += "\r\n"
	msg += bodyText + "\r\n"
	msg += "\r\n"
	msg += "--boundary\r\n"
	msg += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	msg += "\r\n"
	msg += bodyHTML + "\r\n"
	msg += "\r\n"
	msg += "--boundary--\r\n"

	return msg
}

// QueueEmail adds an email to the queue
func (s *EmailService) QueueEmail(userID, organisationID string, to, subject, bodyHTML, bodyText string, scheduledAt time.Time) (*models.EmailQueue, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	orgUUID, err := uuid.Parse(organisationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organisation ID: %w", err)
	}

	queue := &models.EmailQueue{
		UserID:         userUUID,
		OrganisationID: orgUUID,
		To:             to,
		Subject:        subject,
		BodyHTML:       bodyHTML,
		BodyText:       bodyText,
		Status:         models.EmailQueueStatusPending,
		ScheduledAt:    scheduledAt,
		MaxRetries:     3,
	}

	if err := s.db.Create(queue).Error; err != nil {
		return nil, err
	}

	return queue, nil
}

// QueueSequenceEmail queues an email from a sequence step
func (s *EmailService) QueueSequenceEmail(
	userID, organisationID string,
	sequenceID, stepID *string,
	to, subject, bodyHTML, bodyText string,
	delayDays int,
) (*models.EmailQueue, error) {
	scheduledAt := time.Now().AddDate(0, 0, delayDays)

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	orgUUID, err := uuid.Parse(organisationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organisation ID: %w", err)
	}

	queue := &models.EmailQueue{
		UserID:         userUUID,
		OrganisationID: orgUUID,
		To:             to,
		Subject:        subject,
		BodyHTML:       bodyHTML,
		BodyText:       bodyText,
		Status:         models.EmailQueueStatusPending,
		ScheduledAt:    scheduledAt,
		MaxRetries:     3,
	}

	if sequenceID != nil {
		seqUUID, err := uuid.Parse(*sequenceID)
		if err == nil {
			queue.SequenceID = &seqUUID
		}
	}
	if stepID != nil {
		stepUUID, err := uuid.Parse(*stepID)
		if err == nil {
			queue.SequenceStepID = &stepUUID
		}
	}

	if err := s.db.Create(queue).Error; err != nil {
		return nil, err
	}

	return queue, nil
}

// ProcessQueue processes pending emails in the queue
func (s *EmailService) ProcessQueue(batchSize int) (int, int, error) {
	var pending []models.EmailQueue
	
	// Get pending emails that are due
	if err := s.db.Where("status = ? AND scheduled_at <= ?", models.EmailQueueStatusPending, time.Now()).
		Order("scheduled_at ASC").
		Limit(batchSize).
		Find(&pending).Error; err != nil {
		return 0, 0, err
	}

	sent := 0
	failed := 0

	for _, email := range pending {
		// Mark as processing
		s.db.Model(&email).Update("status", models.EmailQueueStatusProcessing)

		err := s.SendEmail(email.To, email.Subject, email.BodyHTML, email.BodyText)
		now := time.Now()

		if err != nil {
			failed++
			email.RetryCount++
			email.ErrorMessage = err.Error()
			email.FailedAt = &now

			if email.RetryCount >= email.MaxRetries {
				email.Status = models.EmailQueueStatusFailed
			} else {
				email.Status = models.EmailQueueStatusPending
				// Reschedule with exponential backoff
				email.ScheduledAt = now.Add(time.Duration(email.RetryCount*email.RetryCount) * time.Hour)
			}

			s.db.Save(&email)
		} else {
			sent++
			email.Status = models.EmailQueueStatusSent
			email.SentAt = &now
			s.db.Save(&email)

			// Log the sent email
			s.logEmail(&email, now)
		}
	}

	return sent, failed, nil
}

// logEmail creates an email log entry
func (s *EmailService) logEmail(queue *models.EmailQueue, sentAt time.Time) {
	log := &models.EmailLog{
		QueueID:        &queue.ID,
		UserID:         queue.UserID,
		OrganisationID: queue.OrganisationID,
		To:             queue.To,
		Subject:        queue.Subject,
		BodyHTML:       queue.BodyHTML,
		BodyText:       queue.BodyText,
		SentAt:         sentAt,
	}

	s.db.Create(log)
}

// RenderTemplate renders an email template with variables
func (s *EmailService) RenderTemplate(templateStr string, vars map[string]interface{}) (string, error) {
	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Helper functions

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
