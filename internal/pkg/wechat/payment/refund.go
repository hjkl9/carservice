package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RefundPayload struct {
	OutTradeNo  string `json:"out_trade_no"`  // 商户订单号
	OutRefundNo string `json:"out_refund_no"` // 商户退款单号
	Reason      string `json:"reason"`        // 退款原因
	NotifyUrl   string `json:"notify_url"`    // 退款结果回调 url
	Amount      struct {
		Refund   int    `json:"refund"`   // 退款金额
		Total    int    `json:"total"`    // 原订单金额
		Currency string `json:"currency"` // 币种: 默认 CNY
	} `json:"amount"`
}

type RefundResponse struct {
	RefundId    string `json:"refund_id"`
	OutRefundNo string `json:"out_refund_no"`
	OutTradeNo  string `json:"out_trade_no"`
	// 枚举值：
	// ORIGINAL：原路退款
	// BALANCE：退回到余额
	// OTHER_BALANCE：原账户异常退到其他余额账户
	// OTHER_BANKCARD：原银行卡异常退到其他银行卡
	// 示例值：ORIGINAL
	Channel string `json:"channel"`
	// 取当前退款单的退款入账方，有以下几种情况：
	// 1）退回银行卡：{银行名称}{卡类型}{卡尾号}
	// 2）退回支付用户零钱:支付用户零钱
	// 3）退还商户:商户基本账户商户结算银行账户
	// 4）退回支付用户零钱通:支付用户零钱通
	// 示例值：招商银行信用卡0403
	UserReceivedAccount string `json:"user_received_account"`
	CreateTime          string `json:"create_time"`
	// 退款到银行发现用户的卡作废或者冻结了，导致原路退款银行卡失败，可前往商户平台-交易中心，手动处理此笔退款。
	// 枚举值：
	// SUCCESS：退款成功
	// CLOSED：退款关闭
	// PROCESSING：退款处理中
	// ABNORMAL：退款异常
	// 示例值：SUCCESS
	Status string `json:"status"`
	// ? need `amount`?
}

func Refund(ctx context.Context, payload RefundPayload) (*RefundResponse, error) {
	api := "https://api.mch.weixin.qq.com/v3/refund/domestic/refunds"

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, api, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response RefundResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	fmt.Printf("status: %s\n", response.Status)
	return &response, nil
}
