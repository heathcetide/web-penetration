package scan

import "errors"

var (
    // 通用错误
    ErrTaskQueueFull    = errors.New("task queue is full")
    ErrActionQueueFull  = errors.New("action queue is full")
    ErrInvalidConfig    = errors.New("invalid configuration")
    ErrNotImplemented   = errors.New("not implemented")
    
    // 扫描错误
    ErrScanTimeout      = errors.New("scan timeout")
    ErrTargetUnreachable = errors.New("target unreachable")
    ErrInvalidTarget    = errors.New("invalid target")
    ErrInvalidPort      = errors.New("invalid port")
    
    // 漏洞检测错误
    ErrRuleNotFound     = errors.New("rule not found")
    ErrInvalidRule      = errors.New("invalid rule")
    ErrDetectionFailed  = errors.New("detection failed")
    
    // 响应错误
    ErrHandlerNotFound  = errors.New("handler not found")
    ErrResponseTimeout  = errors.New("response timeout")
    ErrResponseFailed   = errors.New("response failed")
) 