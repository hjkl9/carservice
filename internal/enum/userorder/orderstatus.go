package userorder

const (
	Pending                 uint8 = 1 // 订单待商家浏览
	AwaitingPayment         uint8 = 2 // 等待支付
	AwaitingAssignInstaller uint8 = 3 // 等待分配 -------------+
	AwaitingInstallation    uint8 = 4 // 等待安装 -- 已支付完成 +
	Completed               uint8 = 5 // 订单完成
	Cancelled               uint8 = 6 // 取消订单
	Refunded                uint8 = 7 // 已退款
)

// OrderStatusDesc 获取字符串订单状态
func OrderStatusDesc(i uint8) string {
	switch i {
	case Pending:
		return "等待商家接单" // 待处理
	case AwaitingPayment:
		return "等待用户付款" // 待确认付款
	case AwaitingAssignInstaller:
		return "等待分配安装师傅"
	case AwaitingInstallation:
		return "待安装"
	case Completed:
		return "已完成"
	case Cancelled:
		return "" // 已取消
	case Refunded:
		return "已退款" // 已退款
	default:
		return "未知状态"
	}
}

const DefaultAtCreation uint8 = Pending
