package v1

type OrderRes struct {
	PayResult interface{} `json:"config"`
	Msg       string      `json:"msg"`
	No        string      `json:"no"`
	PayTypes  string      `json:"pay_types"`
}

type OrderParam struct {
	OrderNo   string `json:"order_no"`
	Coin      string `json:"coin"`
	MID       int64  `json:"m_id"`
	Amount    string `json:"amount"`
	Chain     string `json:"chain"`
	ReturnURL string `json:"return_url"`
	Remark    string `json:"remark"`
}

type OrderListReq struct {
	Page
	OrderNo    string `json:"order_no"`
	MerchantID int64  `json:"merchant_id"`
}

type OrderPayOut struct {
	MerName   string `json:"mer_name"`
	OrderNo   string `json:"order_no"`
	Chain     string `json:"chain"`
	Coin      string `json:"coin"`
	Amount    string `json:"amount"`
	Ac        string `json:"ac"`
	ReturnURL string `json:"return_url"`
}

type OrderPayTxOut struct {
	TxHash string `json:"txHash"`
}
