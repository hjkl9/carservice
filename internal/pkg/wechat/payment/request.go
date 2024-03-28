package payment

// {
//     "transaction_id":"1217752501201407033233368018",
//     "amount":{
//         "payer_total":100,
//         "total":100,
//         "currency":"CNY",
//         "payer_currency":"CNY"
//     },
//     "mchid":"1230000109",
//     "trade_state":"SUCCESS",
//     "bank_type":"CMC",
//     "promotion_detail":[
//         {
//             "amount":100,
//             "wechatpay_contribute":0,
//             "coupon_id":"109519",
//             "scope":"GLOBAL",
//             "merchant_contribute":0,
//             "name":"单品惠-6",
//             "other_contribute":0,
//             "currency":"CNY",
//             "stock_id":"931386",
//             "goods_detail":[
//                 {
//                     "goods_remark":"商品备注信息",
//                     "quantity":1,
//                     "discount_amount":1,
//                     "goods_id":"M1006",
//                     "unit_price":100
//                 },
//                 {
//                     "goods_remark":"商品备注信息",
//                     "quantity":1,
//                     "discount_amount":1,
//                     "goods_id":"M1006",
//                     "unit_price":100
//                 }
//             ]
//         },
//         {
//             "amount":100,
//             "wechatpay_contribute":0,
//             "coupon_id":"109519",
//             "scope":"GLOBAL",
//             "merchant_contribute":0,
//             "name":"单品惠-6",
//             "other_contribute":0,
//             "currency":"CNY",
//             "stock_id":"931386",
//             "goods_detail":[
//                 {
//                     "goods_remark":"商品备注信息",
//                     "quantity":1,
//                     "discount_amount":1,
//                     "goods_id":"M1006",
//                     "unit_price":100
//                 },
//                 {
//                     "goods_remark":"商品备注信息",
//                     "quantity":1,
//                     "discount_amount":1,
//                     "goods_id":"M1006",
//                     "unit_price":100
//                 }
//             ]
//         }
//     ],
//     "success_time":"2018-06-08T10:34:56+08:00",
//     "payer":{
//         "openid":"oUpF8uMuAJO_M2pxb1Q9zNjWeS6o"
//     },
//     "out_trade_no":"1217752501201407033233368018",
//     "AppID":"wxd678efh567hg6787",
//     "trade_state_desc":"支付成功",
//     "trade_type":"MICROPAY",
//     "attach":"自定义数据",
//     "scene_info":{
//         "device_id":"013467007045764"
//     }
// }
type DecryptedResource struct {
	Mchid               string `json:"mchid"`
	TransactionId       string `json:"transactionId"`
	OutTradeNo          string `json:"outTradeNo"`
	RefundId            string `json:"refundId"`
	OutRefundNo         string `json:"outRefundNo"`
	RefundStatus        string `json:"refundStatus"`
	SuccessTime         string `json:"successTime"`
	UserReceivedAccount string `json:"userReceivedAccount"`
	Amount              struct {
		Total       uint64 `json:"total"`
		Refund      uint64 `json:"refund"`
		PayerTotal  uint64 `json:"payerTotal"`
		PayerRefund uint64 `json:"payerRefund"`
	} `json:"amount"`
}
