package idgen

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenNanoId(alphabet string, size int) (id string, err error) {
	if alphabet == "" {
		if size <= 0 {
			return gonanoid.New()
		}
		return gonanoid.New(size)
	}
	return gonanoid.Generate(alphabet, size)
}
