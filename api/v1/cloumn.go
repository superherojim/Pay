package v1

type ColumnParam struct {
	ID          int32  `json:"id"`
	UID         string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Pic         string `json:"pic"`
}

type FileColumnParam struct {
	ID   int32    `json:"id"`
	FIDS []string `json:"fids"`
}

type Columns struct {
	ID          int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	UID         string `gorm:"column:uid" json:"uid"`
	Description string `gorm:"column:description" json:"description"`
	FC          int64  `json:"fc"`
	XID         string `json:"xid"`
	Pic         string `json:"pic"`
}
