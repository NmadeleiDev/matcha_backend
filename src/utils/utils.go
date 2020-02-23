package utils

import (
	"crypto/md5"
	"fmt"
)

func IsEmpty(data string) bool {
	if len(data) == 0 {
		return true
	} else {
		return false
	}
}

func GetMD5(input string) string  {
	h := md5.New()
	result := fmt.Sprintf("%x", h.Sum(nil))
	return result
}