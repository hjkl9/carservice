package userorder

const (
	Pending uint8 = iota
	ToBeConfirmedAndPay
	Cancelled
	Refunded
	Paid
	PrepareToInstall
	InInstallation
	Closed
	Completed
)

// OrderStatusDesc 获取字符串订单状态
func OrderStatusDesc(i uint8) string {
	switch i {
	case Pending:
		return "等待商家接单" // 待处理
	case ToBeConfirmedAndPay:
		return "等待用户确认并付款" // 待确认付款
	case Cancelled:
		return "已取消" // 已取消
	case Refunded:
		return "已退款" // 已退款
	case Paid:
		return "已付款"
	case PrepareToInstall:
		return "待安装"
	case InInstallation:
		return "安装中"
	case Closed:
		return "已关闭"
	case Completed:
		return "已完成"
	default:
		return "未知状态"
	}
}

const DefaultAtCreation uint8 = Pending
