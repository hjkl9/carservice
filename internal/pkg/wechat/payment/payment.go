package payment

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var (
	LoadPrivateKeyErr = errors.New("读取私钥发生错误")
	InitialClientErr  = errors.New("初始化微信支付客户端发生错误")
	PrepayOrderErr    = errors.New("Jsapi 预支付下单调起时发生错误")
)

type PaymentPayload struct {
	Description string // 描述
	OutTradeNo  string // e.g. 订单号码
	Attach      string // 附加数据
	NotifyUrl   string // 回调地址
	Amount      int64  // 交易金额
	OpenId      string // openid
}

type PrepayWithRequestPaymentResponse struct {
	*jsapi.PrepayWithRequestPaymentResponse
}

func PrepayOrder(cfg PaymentConfig, payload PaymentPayload) (*PrepayWithRequestPaymentResponse, error) {
	var (
		mchID               string = cfg.MchId
		mchCertSerialNumber string = cfg.MchCertSerialNumber
		mchAPIv3Key         string = cfg.MchApiV3Key
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	// Given path or string.
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(cfg.PrivateKeyPath)
	if err != nil {
		return nil, LoadPrivateKeyErr
	}
	// mchPrivateKey, err := utils.LoadPrivateKey("YOUR_STRING_PRIVATE_KEY")

	ctx := context.Background()

	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("new wechat pay client err: %s", err)
		return nil, InitialClientErr
	}

	_ = client

	svc := jsapi.JsapiApiService{
		Client: client,
	}

	resp, result, err := svc.PrepayWithRequestPayment(ctx, jsapi.PrepayRequest{
		// Configuration relation.
		Appid: core.String(cfg.Appid),
		Mchid: core.String(cfg.MchId),
		// Payload relation.
		Description: core.String(payload.Description),
		OutTradeNo:  core.String(payload.OutTradeNo),
		Attach:      core.String(payload.Attach),
		NotifyUrl:   core.String(payload.NotifyUrl),
		Amount: &jsapi.Amount{
			Total:    core.Int64(payload.Amount),
			Currency: core.String("CNY"),
		},
	})
	if err != nil {
		return nil, PrepayOrderErr
	}

	// 检查预支付返回的 HTTP 状态码
	if result.Response.StatusCode != http.StatusOK {
		fmt.Printf("请求无效, HTTP Status: %d\n", result.Response.StatusCode)
	}

	// 返回调起预支付的 PrepayID 等等
	_ = resp.PaySign  // 支付签名
	_ = resp.PrepayId // 预支付 ID
	_ = resp.NonceStr // 随机字符串 (非密码)

	return &PrepayWithRequestPaymentResponse{
		resp,
	}, nil
}
