package wechat

import (
	"context"
	"fmt"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

func RequestPayment() error {
	var (
		appId                      string = "UNKNOWN_YOUR_APP_ID"
		privateKeyPath             string = "/path/to/your/merchant/unknown_apiclient_key.pem"
		mchID                      string = "unknown_190000****"                               // 商户号
		mchCertificateSerialNumber string = "unknown_3775B6A45ACD588826D15E583A95F5DD********" // 商户证书序列号
		mchAPIv3Key                string = "unknown_2ab9****************************"         // 商户APIv3密钥
	)
	// load merchant private key file.
	mchPrivateKey, err := utils.LoadPrivateKey(privateKeyPath)
	if err != nil {
		fmt.Println("load merchant private key error, err:" + err.Error())
		return err
	}
	// create a empty context.
	ctx := context.Background()
	// initial client.
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		fmt.Println("create wechat pay client occurs error, err:" + err.Error())
		return err
	}

	// For Native payment.
	svc := native.NativeApiService{Client: client}
	resp, rs, err := svc.Prepay(
		ctx,
		native.PrepayRequest{
			Appid:       core.String(appId), // From config.
			Mchid:       core.String(mchID), // From config.
			Description: core.String("车一一有限公司 - 装车订单"),
			OutTradeNo:  core.String("1217752501201407033233368018"),
			Attach:      core.String("自定义数据说明"),
			NotifyUrl:   core.String("https://www.weixin.qq.com/wxpay/pay.php"),
			Amount: &native.Amount{
				Total: core.Int64(100),
			},
		},
	)
	if err != nil {
		fmt.Println("prepay occurs error, err:" + err.Error())
		return err
	}
	fmt.Printf("status=%d resp=%s\n", rs.Response.StatusCode, resp)
	return nil
}
