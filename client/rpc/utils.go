package rpc

import "encoding/base64"

func tryBase64Decoding(data string) string {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return data
	}
	return string(decoded)
}
