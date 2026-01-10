package driver

import (
	"crypto/sha256"
	"encoding/base64"
)

type CloudInit struct {
	Meta    string
	Network string
	Vendor  string
	User    string
}

func (ci *CloudInit) Hash() string {
	hasher := sha256.New()
	hasher.Write([]byte(ci.Meta))
	hasher.Write([]byte{0})
	hasher.Write([]byte(ci.Network))
	hasher.Write([]byte{0})
	hasher.Write([]byte(ci.Vendor))
	hasher.Write([]byte{0})
	hasher.Write([]byte(ci.User))

	sum := hasher.Sum(nil)

	return base64.RawStdEncoding.EncodeToString(sum)
}
