package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// Test types
type Payload struct {
	Description string // 描述
	OutTradeNo  string // e.g. 订单号码
	Attach      string // 附加数据
	NotifyUrl   string // 回调地址
	Amount      int64  // 交易金额
	OpenId      string // openid
}
type UnifiedOrderFunc func(Payload) error

type PaymentConfig struct {
	MchId               string // 商户号
	MchCertSerialNumber string // 商户证书序列号
	MchApiV3Key         string // 商户 APIv3 密钥
	Appid               string
	PrivateKeyPath      string // 存储私钥路径 `**/**/apiclient_key.pem`
}

type PaymentPayload struct {
	Description string // 描述
	OutTradeNo  string // e.g. 订单号码
	Attach      string // 附加数据
	NotifyUrl   string // 回调地址
	Amount      int64  // 交易金额
	OpenId      string // openid
}

func UnifiedOrder(cfg PaymentConfig, payload PaymentPayload) error {
	var (
		// 商户号 <From configuration>
		mchID string = cfg.MchId
		// 商户证书序列号 <From configuration>
		mchCertificateSerialNumber string = cfg.MchCertSerialNumber
		// 商户 APIv3 密钥 <From configuration>
		mchAPIv3Key string = cfg.MchApiV3Key
	)

	// 获取商户私钥 <From file>
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(cfg.PrivateKeyPath)
	if err != nil {
		// 读取商户私钥失败
		return err
	}

	// 创建空的上下文
	ctx := context.Background()

	// 使用 [商户号, 商户证书序列号, 商户私钥, 商户 APIv3 密钥] 等初始化 client
	// 并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		// 初始化 client 失败
		return err
	}

	svc := jsapi.JsapiApiService{Client: client}
	resp, rs, err := svc.PrepayWithRequestPayment(ctx, jsapi.PrepayRequest{
		Appid:       core.String(cfg.Appid), // Appid from configuration
		Mchid:       core.String(cfg.MchId),
		Description: core.String(payload.Description),
		OutTradeNo:  core.String(payload.OutTradeNo), // 订单号
		Attach:      core.String(payload.Attach),
		NotifyUrl:   core.String("https://www.weixin.qq.com/wxpay/pay.php"), // 回调地址
		Amount: &jsapi.Amount{
			Total: core.Int64(payload.Amount), // 订单金额
		},
		Payer: &jsapi.Payer{
			Openid: core.String(payload.OpenId),
		},
	})
	if err != nil {
		return err
	}
	log.Printf("status=%d resp=%s", rs.Response.StatusCode, resp)

	return nil
}

func PrepayOrder(cfg PaymentConfig, payload PaymentPayload) error {
	var (
		mchID               string = cfg.MchId
		mchCertSerialNumber string = cfg.MchCertSerialNumber
		mchAPIv3Key         string = cfg.MchApiV3Key
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	// Given path or string.
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("/path/to/merchant/apiclient_key.pem")
	if err != nil {
		return err
	}
	// mchPrivateKey, err := utils.LoadPrivateKey("YOUR_STRING_PRIVATE_KEY")

	ctx := context.Background()

	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("new wechat pay client err: %s", err)
		return err
	}

	_ = client

	// TODO

	return nil
}

// UnifiedOrder0110 统一下单
func JsApiPreOrder(cfg PaymentConfig, payload PaymentPayload) error {
	var apiurl = "https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi"

	var data struct {
		AppId       string `json:"appid"`        // 应用 ID
		MchId       string `json:"mchid"`        // 直连商户号 ID
		Description string `json:"description"`  // 商品描述
		OutTradeNo  string `json:"out_trade_no"` // 商户订单号
		Attach      string `json:"attach"`       // 附加数据
		NotifyUrl   string `json:"notify_url"`   // 回调通知地址
		Amount      struct {
			Total    int64  `json:"total"`    // 总金额
			Currency string `json:"currency"` // 货币类型, e.g. CNY
		} `json:"amount"` // 订单金额
		Payer struct {
			OpenId string `json:"openid"` // 用户标识 openid
		} `json:"payer"` // 支付者
	}

	// config
	data.AppId = cfg.Appid
	data.MchId = cfg.MchId
	// payload
	data.Description = payload.Description
	data.OutTradeNo = payload.OutTradeNo
	data.Attach = payload.Attach
	data.NotifyUrl = payload.NotifyUrl
	data.Amount.Total = payload.Amount
	data.Amount.Currency = "CNY"
	data.Payer.OpenId = payload.OpenId

	ctx := context.Background()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, apiurl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// if response.StatusCode != http.StatusOK {
	// 	return errors.New("response status is not ok," + response.Status)
	// }

	rs, _ := io.ReadAll(response.Body)
	fmt.Println(string(rs))
	// reader.

	// decoder := json.NewDecoder(response.Body)
	// if err = decoder.Decode(struct{}{}); err != nil {
	// 	return err
	// }

	return nil
}
