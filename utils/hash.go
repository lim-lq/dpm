package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

func StringToMd5String(text string) string {
	return fmt.Sprintf("%x", StringToMd5Bytes(text))
}

func StringToMd5Bytes(text string) []byte {
	h := md5.New()
	io.WriteString(h, text)
	return h.Sum(nil)
}
