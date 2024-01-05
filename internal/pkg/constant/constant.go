package constant

// Sms related.
var (
	SmsCaptchaPrefix                = "captcha:"
	SmsBindPhoneNumberCaptchaPrefix = "captcha1:"
	SmsNotifyPrefix                 = "notify:"
)

// Websocket related.
var (
	WebsocketServices = map[uint8]string{
		0: "orderNotice",
	}
)
