package constants

import "errors"

var (
	// ErrInvalidCertificates 证书配置错误
	ErrInvalidCertificates = errors.New("certificates must be exactly two")
)
