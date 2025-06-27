package idgen

import "github.com/rs/xid"

func GenXId() string {
	id := xid.New()
	return id.String()
}
