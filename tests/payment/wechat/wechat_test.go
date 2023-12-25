package wechat_test

import (
	"carservice/tests/payment/wechat"
	"fmt"
	"testing"
)

func TestRequestPayment(t *testing.T) {
	err := wechat.RequestPayment()
	if err != nil {
		fmt.Println(err.Error())
	}
}
