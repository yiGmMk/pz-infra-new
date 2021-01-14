package common

type FundBill struct {
	FundChannel string `json:"fundChannel"`
	Amount      string `json:"amount"`
}

// AliWebPayResult 支付宝支付结果回调
type AliWebPayMap struct {
	NotifyTime      string `json:"notify_time"`
	NotifyType      string `json:"notify_type"`
	NotifyID        string `json:"notify_id"`
	SignType        string `json:"sign_type"`
	Sign            string `json:"sign"`
	OutTradeNo      string `json:"out_trade_no"`
	Subject         string `json:"subject"`
	TradeNo         string `json:"trade_no"`
	TradeStatus     string `json:"trade_status"`
	GmtPayMent      string `json:"gmt_payment"`
	SellerEmail     string `json:"seller_email"`
	BuyerLogonId    string `json:"buyer_logon_id"`
	SellerID        string `json:"seller_id"`
	BuyerID         string `json:"buyer_id"`
	TotalAmount     string `json:"total_amount"`
	FundBillList    string `json:"fund_bill_list"`
	MDiscountAmount string `json:"mdiscount_amount"`
	DiscountAmount  string `json:"discount_amount"`
	PassbackParams  string `json:"passback_params"`
}

// AliWebPayResult 支付宝支付结果回调
type AliWebPayResult struct {
	NotifyTime      string     `json:"notify_time"`
	NotifyType      string     `json:"notify_type"`
	NotifyID        string     `json:"notify_id"`
	SignType        string     `json:"sign_type"`
	Sign            string     `json:"sign"`
	OutTradeNo      string     `json:"out_trade_no"`
	Subject         string     `json:"subject"`
	TradeNo         string     `json:"trade_no"`
	TradeStatus     string     `json:"trade_status"`
	GmtPayMent      string     `json:"gmt_payment"`
	SellerEmail     string     `json:"seller_email"`
	BuyerLogonId    string     `json:"buyer_logon_id"`
	SellerID        string     `json:"seller_id"`
	BuyerID         string     `json:"buyer_id"`
	TotalAmount     string     `json:"total_amount"`
	FundBillList    []FundBill `json:"fund_bill_list"`
	MDiscountAmount string     `json:"mdiscount_amount"`
	DiscountAmount  string     `json:"discount_amount"`
	PassbackParams  string     `json:"passback_params"`
}

type AliQueryResult struct {
	TradeNo        string `json:"trade_no"`
	OutTradeNo     string `json:"out_trade_no"`
	OpenID         string `json:"open_id"`
	BuyerLogonID   string `json:"buyer_logon_id"`
	TradeStatus    string `json:"trade_status"`
	TotalAmount    string `json:"total_amount"`
	ReceiptAmount  string `json:"receipt_amount"`
	BuyerPayAmount string `json:"BuyerPayAmount"`
	PointAmount    string `json:"point_amount"`
	InvoiceAmount  string `json:"invoice_amount"`
	SendPayDate    string `json:"send_pay_date"`
	AlipayStoreID  string `json:"alipay_store_id"`
	StoreID        string `json:"store_id"`
	TerminalID     string `json:"terminal_id"`
	FundBillList []struct {
		FundChannel string `json:"fund_channel"`
		Amount      string `json:"amount"`
	} `json:"fund_bill_list"`
	StoreName           string `json:"store_name"`
	BuyerUserID         string `json:"buyer_user_id"`
	DiscountGoodsDetail string `json:"discount_goods_detail"`
	IndustrySepcDetail  string `json:"industry_sepc_detail"`
	PassbackParams      string `json:"passback_params"`
}

type AliRefundResult struct {
	AlipayTradeRefundResp AlipayTradeRefundResponse `json:"alipay_trade_refund_response"`
	Sign                  string                    `json:"sign"`
}

// AliRefundResult 支付宝退款结果
type AlipayTradeRefundResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	//BuyerLogonId string `json:"buyer_logon_id"`
	//BuyerUserId  string `json:"buyer_user_id"`
	//FundChange   string `json:"fund_change"`
	//GmtRefundPay string `json:"gmt_refund_pay"`
	//TradeNo      string `json:"tran_no"`
	//OutTradeNo   string `json:"out_tran_no"`
	//RefundFee    string `json:"refund_fee"`
	//SendBackFee  string `json:"send_back_fee"`
}
