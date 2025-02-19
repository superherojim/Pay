package xid

import (
	"github.com/rs/xid"
)

func GenXID() string {
	return xid.New().String()
}
