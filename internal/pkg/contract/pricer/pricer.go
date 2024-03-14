package pricer

type FloatPrice = float64
type UintPrice = uint64

// 物品价格
type ItemPricer interface {
	ItemId() uint64
	GetUint64() uint64
	GetFloat64() float64
}

func ComputeTotal(items []ItemPricer) (FloatPrice, UintPrice) {
	var i UintPrice = 0
	var f FloatPrice = 0.00

	for _, item := range items {
		i += item.GetUint64()
		f += item.GetFloat64()
	}

	return f, i
}

func Test(item ItemPricer) {
	return
}
