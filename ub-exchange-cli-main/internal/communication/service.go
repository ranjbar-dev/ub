package communication

import (
	"bytes"
	"encoding/json"
	"exchange-go/internal/platform"
	"html/template"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	TemplatesPathPrefix = "./internal/communication/templates/"
	TypeEmail           = "email"
	TypeSMS             = "sms"

	PaymentTypeWithdraw = "WITHDRAW"
	PaymentTypeDeposit  = "DEPOSIT"

	PaymentStatusCreated    = "CREATED"
	PaymentStatusCompleted  = "COMPLETED"
	PaymentStatusInProgress = "IN_PROGRESS"
	PaymentStatusFailed     = "FAILED"
	PaymentStatusCanceled   = "CANCELED"
	PaymentStatusRejected   = "REJECTED"
)

type CommunicatingUser struct {
	Email string
	Phone string
}

type publishData struct {
	Receiver    string `json:"receiver"`
	Subject     string `json:"subject"`
	Content     string `json:"content"`
	Priority    int    `json:"priority"`
	Type        string `json:"type"`
	ScheduledAt string `json:"scheduledAt"`
}

type CryptoPaymentStatusUpdateEmailParams struct {
	FullName        string
	Type            string
	Status          string
	CurrencyCode    string
	Amount          string
	ToAddress       string
	RejectionReason string
}

type ContactUsToAdminParams struct {
	Name    string
	Email   string
	Subject string
	Body    string
}

type PasswordChangedEmailParams struct {
	Email         string
	CurrentIP     string
	CurrentDevice string
	ChangedDate   string
}

type ForgotPasswordEmailParams struct {
	Link          string
	CurrentIP     string
	CurrentDevice string
	RequestDate   string
}

type SendIPChangedEmailParams struct {
	Email       string
	FullName    string
	LastIP      string
	CurrentIP   string
	Device      string
	ChangedDate string
}

// Service is the transactional email and SMS communication service. It renders message
// templates and publishes them to the RabbitMQ queue for asynchronous delivery.
type Service interface {
	// SendWithdrawConfirmationEmail sends a withdrawal confirmation email containing the
	// confirmation code, currency, amount, and destination address.
	SendWithdrawConfirmationEmail(u CommunicatingUser, coin, amount, address, code string)
	// SendVerificationEmailToUser sends an email verification link to a newly registered user.
	SendVerificationEmailToUser(u CommunicatingUser, link string)
	// SendCryptoPaymentStatusUpdateEmail sends an email notifying the user about a crypto
	// payment status change (e.g., deposit completed, withdrawal rejected).
	SendCryptoPaymentStatusUpdateEmail(u CommunicatingUser, params CryptoPaymentStatusUpdateEmailParams)
	// SendContactUsToAdmin forwards a contact-us message from a user to the admin email address.
	SendContactUsToAdmin(adminUser CommunicatingUser, params ContactUsToAdminParams)
	// SendPasswordChangedEmail sends a notification email when a user's password has been changed.
	SendPasswordChangedEmail(u CommunicatingUser, params PasswordChangedEmailParams)
	// SendUserPhoneConfirmationSms sends an SMS with a confirmation code for phone verification.
	SendUserPhoneConfirmationSms(u CommunicatingUser, code string)
	// SendUserForgotPasswordEmail sends a forgot-password email with a reset link.
	SendUserForgotPasswordEmail(u CommunicatingUser, params ForgotPasswordEmailParams)
	// SendIPChangedEmail sends an alert email when a user's login IP address has changed.
	SendIPChangedEmail(u CommunicatingUser, params SendIPChangedEmailParams)
}

type service struct {
	queueManager QueueManager
	logger       platform.Logger
}

func (s *service) SendVerificationEmailToUser(u CommunicatingUser, link string) {
	type context struct {
		Link string
	}

	data := context{
		Link: link,
	}
	receivers := []CommunicatingUser{u}
	s.sendMessage(TemplateCodeUserEmailVerification, "", receivers, data, TypeEmail, 5, time.Now())

}

func (s *service) SendWithdrawConfirmationEmail(u CommunicatingUser, coin, amount, address, code string) {
	type context struct {
		Code         string
		CurrencyCode string
		Amount       string
		ToAddress    string
	}

	data := context{
		Code:         code,
		CurrencyCode: strings.ToUpper(coin),
		Amount:       amount,
		ToAddress:    address,
	}
	receivers := []CommunicatingUser{u}
	s.sendMessage(TemplateCryptoPaymentWithdrawConfirmationEmail, "", receivers, data, TypeEmail, 5, time.Now())
}

func (s *service) SendCryptoPaymentStatusUpdateEmail(u CommunicatingUser, params CryptoPaymentStatusUpdateEmailParams) {
	type context struct {
		FullName        string
		CurrencyCode    string
		Amount          string
		ToAddress       string
		RejectionReason string
	}

	templateName := ""

	if params.Type == PaymentTypeDeposit {
		switch params.Status {
		case PaymentStatusCompleted:
			templateName = TemplateCryptoPaymentDepositCompletedEmail
			break
		case PaymentStatusFailed:
			//templateName = TemplateCryptoPaymentDepositFailedEmail
			templateName = ""
			break
		default:
			templateName = ""
		}
	} else {
		switch params.Status {
		case PaymentStatusCreated:
			templateName = TemplateCryptoPaymentWithdrawCreatedEmail
			break
		case PaymentStatusCompleted:
			//templateName = TemplateCryptoPaymentWithdrawCompletedEmail
			templateName = ""
			break
		case PaymentStatusFailed:
			templateName = TemplateCryptoPaymentWithdrawFailedEmail
			break
		case PaymentStatusRejected:
			templateName = TemplateCryptoPaymentWithdrawRejectedEmail
			break
		default:
			templateName = ""
		}

	}

	if templateName != "" {
		data := context{
			FullName:        params.FullName,
			CurrencyCode:    params.CurrencyCode,
			Amount:          params.Amount,
			ToAddress:       params.ToAddress,
			RejectionReason: params.RejectionReason,
		}

		receivers := []CommunicatingUser{u}
		s.sendMessage(templateName, "", receivers, data, TypeEmail, 5, time.Now())
	}
}

func (s *service) sendMessage(templateName string, messageTitle string, receivers []CommunicatingUser, context interface{}, messageType string, priority int, scheduledAt time.Time) {
	if messageTitle == "" {
		messageTitle, _ = templateTitles[templateName]
	}

	messageContent, err := s.getContent(templateName, context)
	if err != nil {
		s.logger.Error2("can not render template", err,
			zap.String("service", "communicationService"),
			zap.String("method", "sendMessage"),
			zap.String("templateName", templateName),
		)
		return
	}
	for _, u := range receivers {
		data := publishData{
			Subject:     messageTitle,
			Content:     messageContent,
			Priority:    priority,
			Type:        messageType,
			ScheduledAt: scheduledAt.Format("2006-01-02 15:04:05"),
		}
		if messageType == TypeEmail {
			data.Receiver = u.Email
			s.publishEmail(data)
		}

		if messageType == TypeSMS {
			data.Receiver = u.Phone
			s.publishSMS(data)
		}
	}
}

func (s *service) getContent(templateName string, context interface{}) (string, error) {
	templateCompleteName := templateName + ".html"
	t := template.New(templateCompleteName)
	file := TemplatesPathPrefix + templateCompleteName
	t, err := t.ParseFiles(file)
	if err != nil {
		return "", err
	}

	var writer bytes.Buffer
	err = t.Execute(&writer, context)
	if err != nil {
		return "", err
	}

	return writer.String(), nil
}

func (s *service) SendContactUsToAdmin(adminUser CommunicatingUser, params ContactUsToAdminParams) {
	type context struct {
		Name    string
		Email   string
		Subject string
		Body    string
	}

	data := context{
		Name:    params.Name,
		Email:   params.Email,
		Subject: params.Subject,
		Body:    params.Body,
	}
	receivers := []CommunicatingUser{adminUser}
	s.sendMessage(TemplateContactUsMessage, "", receivers, data, TypeEmail, 5, time.Now())
}

func (s *service) SendPasswordChangedEmail(u CommunicatingUser, params PasswordChangedEmailParams) {
	type context struct {
		Email         string
		CurrentIP     string
		CurrentDevice string
		ChangedDate   string
	}

	data := context{
		Email:         params.Email,
		CurrentIP:     params.CurrentIP,
		CurrentDevice: params.CurrentDevice,
		ChangedDate:   params.ChangedDate,
	}

	receivers := []CommunicatingUser{u}
	s.sendMessage(TemplatePasswordHasBeenChanged, "", receivers, data, TypeEmail, 5, time.Now())
}

func (s *service) SendUserPhoneConfirmationSms(u CommunicatingUser, code string) {
	type context struct {
		Code string
	}

	data := context{
		Code: code,
	}

	receivers := []CommunicatingUser{u}
	s.sendMessage(TemplatePhoneConfirmationSms, "", receivers, data, TypeSMS, 5, time.Now())
}

func (s *service) SendUserForgotPasswordEmail(u CommunicatingUser, params ForgotPasswordEmailParams) {
	type context struct {
		Link          string
		Email         string
		CurrentIP     string
		CurrentDevice string
		RequestDate   string
	}

	data := context{
		Link:          params.Link,
		Email:         u.Email,
		CurrentIP:     params.CurrentIP,
		CurrentDevice: params.CurrentDevice,
		RequestDate:   params.RequestDate,
	}

	receivers := []CommunicatingUser{u}
	s.sendMessage(TemplateAuthForgotPassword, "", receivers, data, TypeEmail, 5, time.Now())
}

func (s *service) SendIPChangedEmail(u CommunicatingUser, params SendIPChangedEmailParams) {
	type context struct {
		Email       string
		FullName    string
		LastIP      string
		CurrentIP   string
		Device      string
		ChangedDate string
	}

	data := context{
		Email:       params.Email,
		FullName:    params.FullName,
		LastIP:      params.LastIP,
		CurrentIP:   params.CurrentIP,
		Device:      params.Device,
		ChangedDate: params.ChangedDate,
	}

	receivers := []CommunicatingUser{u}
	title := "[Admin] " + templateTitles["TemplateUserLoginIpChangedEmail"]
	s.sendMessage(TemplateUserLoginIPChangedEmail, title, receivers, data, TypeEmail, 5, time.Now())
}

func (s *service) publishEmail(data publishData) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		s.logger.Error2("can not marshal data", err,
			zap.String("service", "communicationService"),
			zap.String("method", "publishEmail"),
		)
		return
	}
	s.queueManager.PublishEmailOrSms(dataByte)
}

func (s *service) publishSMS(data publishData) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		s.logger.Error2("can not marshal data", err,
			zap.String("service", "communicationService"),
			zap.String("method", "publishSMS"),
		)
		return
	}
	s.queueManager.PublishEmailOrSms(dataByte)
}

func NewCommunicationService(queueManager QueueManager, logger platform.Logger) Service {
	return &service{
		queueManager: queueManager,
		logger:       logger,
	}
}
