package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/1/17 20:33
 * @file: gen_sign.go
 * @description: lark sign gen
 */

func GenSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	strToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret

	var b []byte
	h := hmac.New(sha256.New, []byte(strToSign))
	_, err := h.Write(b)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, err
}
