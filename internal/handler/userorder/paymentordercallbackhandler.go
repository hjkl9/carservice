package userorder

import (
	"crypto/x509"
	"net/http"
	"time"

	"carservice/internal/logic/userorder"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/httper/api"
	"carservice/internal/pkg/wechat/payment"
	"carservice/internal/svc"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"github.com/zeromicro/go-zero/core/logc"
)

func PaymentOrderCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := userorder.NewPaymentOrderCallbackLogic(r.Context(), svcCtx)

		certPath := "need to add in config file"
		// cert, err := utils.LoadCertificate(certString)
		cert, err := utils.LoadCertificateWithPath(certPath)
		if err != nil {
			logc.Errorf(r.Context(), "加载证书时发生错误: %#s\n", err)
			api.ResponseWithCtx(r.Context(), w, nil, errcode.SystemErr.SetMessage(err.Error()))
			return
		}

		handler, err := notify.NewRSANotifyHandler(
			svcCtx.Config.WechatPayMerchantConf.MchApiV3Key,
			verifiers.NewSHA256WithRSAVerifier(core.NewCertificateMapWithList([]*x509.Certificate{cert})),
		)
		if err != nil {
			logc.Errorf(r.Context(), "创建通知处理器时发生错误: %#s\n", err)
			api.ResponseWithCtx(r.Context(), w, nil, errcode.SystemErr.SetMessage(err.Error()))
		}

		type contentType struct {
			Mchid           *string    `json:"mchid"`
			Appid           *string    `json:"appid"`
			CreateTime      *time.Time `json:"create_time"`
			OutContractCode *string    `json:"out_contract_code"`
		}

		decryptedContent := new(payment.DecryptedResource)

		_, err = handler.ParseNotifyRequest(r.Context(), r, decryptedContent)
		if err != nil {
			logc.Errorf(r.Context(), "解析请求时发生错误: %#s\n", err)
			api.ResponseWithCtx(r.Context(), w, nil, errcode.InvalidParametersErr.SetMessage(err.Error()))
			return
		}

		err = l.PaymentOrderCallback(decryptedContent)
		api.ResponseWithCtx(r.Context(), w, nil, err)
	}
}

func callbackMiddleware() {
	// todo
}
