package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// EncodeWithSHA256 加密图片名称
func EncodeWithSHA256(name string) string {
	hash := sha256.Sum256([]byte(name))
	return hex.EncodeToString(hash[:])
}
