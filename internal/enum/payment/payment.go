package payment

type PaymentMethodType uint8

const (
	Unknown PaymentMethodType = iota
	Alipay
	WechatPay
	UnionPay
)

func PaymentMethodDesc(i PaymentMethodType) string {
	switch i {
	case Unknown:
		return "未知"
	case Alipay:
		return "支付宝"
	case WechatPay:
		return "微信"
	case UnionPay:
		return "银联"
	default:
		return "未知支付方式"
	}
}

const DefaultAtCreation PaymentMethodType = Unknown
