package idgen

import (
	"errors"
	"github.com/google/uuid"
)

// GenUUId 根据版本生成 UUID。
// ver: 支持 "v1", "v3", "v4", "v5"（不支持 v2）
// name 和 namespace 仅在 v3 / v5 时需要。
func GenUUId(ver string, name string) (string, error) {
	switch ver {
	case "v1":
		id, err := uuid.NewUUID()
		if err != nil {
			return "", err
		}
		return id.String(), nil

	case "v3":
		if name == "" {
			return "", errors.New("v3版本需要提供name参数")
		}
		return uuid.NewMD5(uuid.NameSpaceDNS, []byte(name)).String(), nil

	case "v4":
		return uuid.New().String(), nil

	case "v5":
		if name == "" {
			return "", errors.New("v5版本需要提供name参数")
		}
		return uuid.NewSHA1(uuid.NameSpaceDNS, []byte(name)).String(), nil

	default:
		return "", errors.New("不支持的UUID版本: " + ver)
	}
}
