package payment

type PaymentConfig struct {
	MchId               string // 商户号
	MchCertSerialNumber string // 商户证书序列号
	MchApiV3Key         string // 商户 APIv3 密钥
	Appid               string
	PrivateKeyPath      string // 存储私钥路径 `**/**/apiclient_key.pem`
}
