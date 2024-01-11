package payment_test

import (
	"carservice/internal/pkg/wechat/payment"
	"testing"
)

func TestPayment(t *testing.T) {
	c := payment.PaymentConfig{
		MchId:          "1900006891",
		PrivateKeyPath: "P:\\payment\\apiclient_key.pem",
	}
	p := payment.PaymentPayload{
		Description: "红内裤",
		OutTradeNo:  "O1010100101",
		Attach:      "穿着红色内裤的人在打电话",
		NotifyUrl:   "https://...",
		Amount:      1,
		OpenId:      "oLBnj5KngUW1T_Es3wTlynHwi-4g", // Mock openid.
	}
	// if err := payment.JsApiPreOrder(c, p); err != nil {
	// 	t.Fatalf("payment failed, err: %s\n", err.Error())
	// }

	err := payment.PrepayOrder(c, p)
	if err != nil {
		t.Fatalf("prepay order failed, err: %s\n", err.Error())
	}
}
