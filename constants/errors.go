package constants

import "errors"

var (
	// ErrMsgTooLager 消息过长
	ErrMsgTooLager = errors.New("message too large")
	// ErrInvalidCertificates 证书配置错误
	ErrInvalidCertificates = errors.New("certificates must be exactly two")
	// ErrNodeStopped 节点已停止
	ErrNodeStopped = errors.New("node: the node has been stopped")
	// ErrMetricNotKnown 未知的指标
	ErrMetricNotKnown = errors.New("the provided metric does not exist")
)
