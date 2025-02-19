package v1

type MerApiParam struct {
	ID          int64  `json:"id"`
	MID         int64  `json:"m_id"`
	CallbackURL string `json:"callback_url"`
	SecretKey   string `json:"secret_key"`
	Remark      string `json:"remark"`
}

type MerApiListReq struct {
	Page
	MID int64 `json:"mid"`
}
