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
	// ErrReporterClosed 指标回报关闭
	ErrReporterClosed = errors.New("metric reporter closed")

	// ErrSessionClosed session已关闭
	ErrSessionClosed = errors.New("session closed")
	// ErrorPushBufferFull session push缓冲已满
	ErrorPushBufferFull = errors.New("session push full")

	// ErrorMsgDiscard 消息被丢弃
	ErrorMsgDiscard = errors.New("message discard")
)
