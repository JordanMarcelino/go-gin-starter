package encryptutils

import "encoding/base64"

type Base64Encryptor interface {
	EncodeURL(src string) string
	EncodeStd(src string) string
	DecodeURL(s string) (string, error)
	DecodeStd(s string) (string, error)
}

type base64Encryptor struct {
}

func NewBase64Encryptor() *base64Encryptor {
	return &base64Encryptor{}
}

func (e *base64Encryptor) EncodeURL(src string) string {
	return base64.URLEncoding.EncodeToString([]byte(src))
}

func (e *base64Encryptor) EncodeStd(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

func (e *base64Encryptor) DecodeURL(s string) (string, error) {
	decodedBytes, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

func (e *base64Encryptor) DecodeStd(s string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return string(decodedBytes), nil
}
