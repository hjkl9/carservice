package userorder

// Should use type of int8.
type OrderStatusType int

const (
	Pending OrderStatusType = iota
	ToBeAcceptedByUser
	ToBePaid
	Cancelled
	Refunded
	Paid
	PrepareToInstall
	InInstallation
	Closed
	Completed
)

// OrderStatusDesc 获取字符串订单状态
func OrderStatusDesc(i OrderStatusType) string {
	switch i {
	case Pending:
		return "等待商家接单"
	case ToBeAcceptedByUser:
		return "等待用户接受报价"
	case ToBePaid:
		return "等待用户付款"
	case Cancelled:
		return "用户取消订单"
	case Refunded:
		return "已退款"
	case Paid:
		return "已付款"
	case PrepareToInstall:
		return "待安装"
	case InInstallation:
		return "安装中"
	case Closed:
		return "订单已关闭"
	case Completed:
		return "订单完成"
	default:
		return "未知状态"
	}
}
