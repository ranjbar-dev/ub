package communication

const (
	TemplateCodeUserEmailVerification                = "user_email_verification"
	TemplateUserProfileStatusIncomplete              = "user_profile_status_incomplete"
	TemplateUserProfileStatusProcessing              = "user_profile_status_processing"
	TemplateUserProfileStatusConfirmed               = "user_profile_status_confirmed"
	TemplateUserProfileStatusRejected                = "user_profile_status_rejected"
	TemplateUserLoginIPChangedEmail                  = "user_login_ip_changed_email"
	TemplateReceiveUserProfile                       = "receive_user_profile"
	TemplateChangeUserLevel                          = "change_user_level"
	TemplateAuthForgotPassword                       = "auth_forgot_password"
	TemplateAuthForgotPasswordForMobile              = "auth_forgot_password_for_mobile"
	TemplatePasswordHasBeenChanged                   = "password_has_been_changed"
	TemplatePhoneConfirmationSms                     = "phone_confirmation_sms"
	TemplateCryptoPaymentDepositCompletedEmail       = "crypto_payment_deposit_completed_email"
	TemplateCryptoPaymentDepositFailedEmail          = "crypto_payment_deposit_failed_email"
	TemplateCryptoPaymentWithdrawCreatedEmail        = "crypto_payment_withdraw_created_email"
	TemplateCryptoPaymentWithdrawCompletedEmail      = "crypto_payment_withdraw_completed_email"
	TemplateCryptoPaymentWithdrawFailedEmail         = "crypto_payment_withdraw_failed_email"
	TemplateCryptoPaymentWithdrawRejectedEmail       = "crypto_payment_withdraw_rejected_email"
	TemplateCryptoPaymentWithdrawConfirmationEmail   = "crypto_payment_withdraw_confirmation_email"
	TemplateLoginEmail                               = "login_email"
	TemplateContactUsMessage                         = "contact_us_message"
	TemplateUserProfileImageStatusRejected           = "user_profile_image_status_rejected"
	TemplateUserProfileImageStatusConfirmed          = "user_profile_image_status_confirmed"
	TemplateUserProfileImageStatusPartiallyConfirmed = "user_profile_image_status_partially_confirmed"
	TemplateUserDeviceConfirmation                   = "user_device_confirmation"
)

var templateTitles = map[string]string{
	TemplateCodeUserEmailVerification:                "Email verification",
	TemplateUserProfileStatusIncomplete:              "Profile data is incomplete",
	TemplateUserProfileStatusProcessing:              "Profile data is in progress",
	TemplateUserProfileStatusConfirmed:               "Profile data has been confirmed",
	TemplateUserProfileStatusRejected:                "Profile data has been rejected",
	TemplateUserLoginIPChangedEmail:                  "Login IP has been changed",
	TemplateReceiveUserProfile:                       "User profile data has been received",
	TemplateChangeUserLevel:                          "Your user level has been changed",
	TemplateAuthForgotPassword:                       "Reset Password",
	TemplatePasswordHasBeenChanged:                   "Password changed",
	TemplateCryptoPaymentDepositCompletedEmail:       "New Deposit has been completed",
	TemplateCryptoPaymentDepositFailedEmail:          "New Deposit has been failed",
	TemplateCryptoPaymentWithdrawCreatedEmail:        "Withdrawal request has been submitted",
	TemplateCryptoPaymentWithdrawCompletedEmail:      "Withdrawal request has been completed",
	TemplateCryptoPaymentWithdrawFailedEmail:         "Withdrawal request has been failed",
	TemplateCryptoPaymentWithdrawRejectedEmail:       "Withdrawal request has been rejected",
	TemplateCryptoPaymentWithdrawConfirmationEmail:   "New withdraw request confirmation",
	TemplateUserProfileImageStatusRejected:           "Identity image has been rejected",
	TemplateUserProfileImageStatusPartiallyConfirmed: "Identity image has been partially confirmed",
	TemplateUserProfileImageStatusConfirmed:          "Identity image has been confirmed",
	TemplateUserDeviceConfirmation:                   "New device login",
	TemplateLoginEmail:                               "Login",
	TemplateContactUsMessage:                         "New contact us message",
}
