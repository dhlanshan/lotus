package idgen

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func GenULId() string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)

	return id.String()
}
