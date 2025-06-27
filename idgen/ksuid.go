package idgen

import "github.com/segmentio/ksuid"

func GenKsuId() string {
	id := ksuid.New()

	return id.String()
}
