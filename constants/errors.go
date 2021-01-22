package constants

import "errors"

var (
	// ErrMsgTooLager 消息过长
	ErrMsgTooLager = errors.New("message too large")
	// ErrInvalidCertificates 证书配置错误
	ErrInvalidCertificates = errors.New("certificates must be exactly two")
)
