package payment_test

import (
	"carservice/internal/pkg/wechat/payment"
	"context"
	"fmt"
	"testing"
)

func TestRefund(t *testing.T) {
	ctx := context.Background()

	resp, err := payment.Refund(ctx, payment.RefundPayload{
		OutTradeNo:  "1111111111",
		OutRefundNo: "222222222222",
		Reason:      "不知道",
		NotifyUrl:   "https://localhost:8888",
		Amount: struct {
			Refund   int    `json:"refund"`
			Total    int    `json:"total"`
			Currency string `json:"currency"`
		}{
			Refund:   1,
			Total:    1,
			Currency: "CNY",
		},
	})
	if err != nil {
		t.Errorf("failed to request refund, err: %s\n", err.Error())
		t.Fail()
	}

	fmt.Printf("%+v\n", resp)
}
