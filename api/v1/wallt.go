package v1

type Wallet struct {
	ID     int64  `json:"id"`
	MID    int64  `json:"m_id"`
	Remark string `json:"remark"` // 备注
}

type AddWallet struct {
	ID       int64  `json:"id"`
	MID      int64  `json:"m_id"`
	Ac       string `json:"ac"`       // 钱包地址
	PriKey   string `json:"pri_key"`  // 私钥
	Mnemonic string `json:"mnemonic"` // 助词器
	Path     string `json:"path"`     // 钱包路径
	Remark   string `json:"remark"`
}

type WalletList struct {
	MID int64 `json:"m_id"`
	Page
}
