// Code generated by goctl. DO NOT EDIT.
package types

type SendCaptchaReq struct {
	PhoneNumber string `form:"phoneNumber"`
}

type SendCaptchaRep struct {
}

type PhoneNumberLoginReq struct {
	PhoneNumber string `json:"phoneNumber"`
	Captcha     string `json:"captcha"`
}

type PhoneNumberLoginRep struct {
	Token string `json:"token"`
}
