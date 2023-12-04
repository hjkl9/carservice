package sms

import (
	"carservice/internal/config"
	stderrors "errors"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

const endpoint = "ap-guangzhou"

type Sms struct {
	config config.Config
}

func NewSms(c config.Config) *Sms {
	return &Sms{
		c,
	}
}

func (s *Sms) newCredential() *common.Credential {
	return common.NewCredential(s.config.SmsConf.SecretId, s.config.SmsConf.SecretKey)
}

func (s *Sms) newClientProfile() *profile.ClientProfile {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	return cpf
}

func (s *Sms) newClient(cr *common.Credential, cpf *profile.ClientProfile) (*sms.Client, error) {
	return sms.NewClient(cr, endpoint, cpf)
}

func (s *Sms) Send(templateIdSet []string, templateSet []string, phoneNumberSet []string) error {
	credential := s.newCredential()
	cpf := s.newClientProfile()
	client, err := s.newClient(credential, cpf)
	if err != nil {
		return err
	}

	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(s.config.SmsConf.SdkAppId)
	request.SignName = common.StringPtr(s.config.SmsConf.SignName)
	if len(templateIdSet) > 0 {
		request.TemplateId = common.StringPtr(templateIdSet[0])
	} else {
		request.TemplateId = common.StringPtr(s.config.SmsConf.TemplateId)
	}
	request.TemplateParamSet = common.StringPtrs(templateSet)
	request.PhoneNumberSet = common.StringPtrs(phoneNumberSet)

	// First variable name is `resp`
	resp, err := client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		// logger.Debugw("TencentCloudSDKError", err.Error())
		return err
	}
	// ? is a `insufficient balance in SMS package` an error?
	if err = s.checkSendSmsStatus(resp); err != nil {
		return err
	}
	// 非 SDK 异常而直接失败
	// 当出现以下错误码时，快速解决方案参考
	// [FailedOperation.SignatureIncorrectOrUnapproved](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Afailedoperation.signatureincorrectorunapproved-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	// [FailedOperation.TemplateIncorrectOrUnapproved](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Afailedoperation.templateincorrectorunapproved-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	// [UnauthorizedOperation.SmsSdkAppIdVerifyFail](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Aunauthorizedoperation.smssdkappidverifyfail-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	// [UnsupportedOperation.ContainDomesticAndInternationalPhoneNumber](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Aunsupportedoperation.containdomesticandinternationalphonenumber-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	// 更多错误，可咨询[腾讯云助手](https://tccc.qcloud.com/web/im/index.html#/chat?webAppId=8fa15978f85cb41f7e2ea36920cb3ae1&title=Sms)
	if err != nil {
		// logger.Debugw("TencentCloudSDKError", err.Error())
		return err
	}
	// if len(resp.Response.SendStatusSet) > 0 {
	// 	if *resp.Response.SendStatusSet[0] == sms.FAILEDOPERATION_INSUFFICIENTBALANCEINSMSPACKAGE {

	// 	}
	// }

	return nil
}

func (s *Sms) checkSendSmsStatus(resp *sms.SendSmsResponse) error {
	statusSet := resp.Response.SendStatusSet
	// get first.
	if len(statusSet) > 0 {
		return stderrors.New(*statusSet[0].Message)
	}
	return nil
}
