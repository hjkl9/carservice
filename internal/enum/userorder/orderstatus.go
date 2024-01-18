package userorder

type OrderStatus struct {
	Status string `json:"status"`
	Label  string `json:"label"`
	GoTag  string `json:"goTag"`
}

func NewOrderStatus(status, label, goTag string) *OrderStatus {
	return &OrderStatus{status, label, goTag}
}

const (
	Pending                 uint8 = 0 // 订单待商家浏览
	AwaitingPayment         uint8 = 1 // 等待支付
	AwaitingAssignInstaller uint8 = 2 // 等待分配 -------------+
	AwaitingInstallation    uint8 = 3 // 等待安装 -- 已支付完成 +
	Completed               uint8 = 4 // 订单完成
	Cancelled               uint8 = 5 // 取消订单
	Refunded                uint8 = 6 // 已退款
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

func ClientTabList() (uint8, [6]*OrderStatus) {
	return 6, [6]*OrderStatus{
		NewOrderStatus("0", "待处理", "@deprecated"),
		NewOrderStatus("1", "待支付", "@deprecated"),
		NewOrderStatus("2", "待安装", "@deprecated"),
		NewOrderStatus("3", "已完成", "@deprecated"),
		NewOrderStatus("4", "已取消", "@deprecated"),
		NewOrderStatus("5", "已退款", "@deprecated"),
	}
}

const DefaultAtCreation uint8 = Pending
