package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

//
// MakeSha1
// @Description: 对密码进行加密
// @param password
// @return string
//
func MakeSha1(password string) string {
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(password))
	return hex.EncodeToString(sha1Hash.Sum(nil))
}
