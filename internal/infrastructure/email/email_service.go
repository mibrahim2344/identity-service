package email

import (
	"fmt"

	"github.com/mibrahim2344/identity-service/internal/domain/services"
)

type emailService struct {
	// Add configuration fields here (SMTP settings, templates, etc.)
}

// NewEmailService creates a new instance of the email service
func NewEmailService() services.EmailService {
	return &emailService{}
}

func (s *emailService) SendEmail(to string, subject string, body string) error {
	// TODO: Implement actual email sending logic
	return fmt.Errorf("email sending not implemented")
}

func (s *emailService) SendVerificationEmail(to string, verificationToken string) error {
	subject := "Verify Your Email"
	body := fmt.Sprintf("Please verify your email using this token: %s", verificationToken)
	return s.SendEmail(to, subject, body)
}

func (s *emailService) SendPasswordResetEmail(to string, resetToken string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf("Use this token to reset your password: %s", resetToken)
	return s.SendEmail(to, subject, body)
}
