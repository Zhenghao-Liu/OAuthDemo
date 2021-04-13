package utils

import (
	"encoding/base64"
)

func Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func Decode(data string) string {
	ans, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ""
	}
	return string(ans)
}
