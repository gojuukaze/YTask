package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func GetStrMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	md5Data := h.Sum([]byte(""))
	return hex.EncodeToString(md5Data)
}

func GetStrSha1(data string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(data))
	return hex.EncodeToString(sha1.Sum([]byte("")))
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
