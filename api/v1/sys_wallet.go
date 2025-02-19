package v1

type SysWalletCreateReq struct {
	Ac       string `json:"ac" binding:"required"`
	PriKey   string `json:"pri_key" binding:"required"`
	Mnemonic string `json:"mnemonic" binding:"required"`
	Path     string `json:"path" binding:"required"`
	Remark   string `json:"remark"`
}

type SysWalletUpdateReq struct {
	ID       int64  `json:"id" binding:"required"`
	Ac       string `json:"ac"`
	PriKey   string `json:"pri_key"`
	Mnemonic string `json:"mnemonic"`
	Path     string `json:"path"`
	Remark   string `json:"remark"`
}
