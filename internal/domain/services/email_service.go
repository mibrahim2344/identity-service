package services

// EmailService defines the interface for sending emails
type EmailService interface {
	// SendEmail sends an email with the given parameters
	SendEmail(to string, subject string, body string) error
	
	// SendVerificationEmail sends a verification email to the user
	SendVerificationEmail(to string, verificationToken string) error
	
	// SendPasswordResetEmail sends a password reset email to the user
	SendPasswordResetEmail(to string, resetToken string) error
}
