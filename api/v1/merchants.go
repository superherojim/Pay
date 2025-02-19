package v1

type MerParam struct {
	ID           int64  `json:"id"`
	Nickname     string `json:"nickname"`     // 昵称
	Avatar       string `json:"avatar"`       // 头像
	Introduction string `json:"introduction"` // 简介
	Password     string `json:"password"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	CreatedBy    string `json:"created_by"`
}

type MerIN struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

type MerchantsListReq struct {
	Page
}
