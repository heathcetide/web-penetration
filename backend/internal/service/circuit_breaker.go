package service

import (
    "errors"
    "sync"
    "time"
)

// 熔断器状态
type CircuitState int

const (
    StateClosed CircuitState = iota    // 关闭状态(正常)
    StateOpen                          // 打开状态(熔断)
    StateHalfOpen                      // 半开状态(尝试恢复)
)

// 熔断器配置
type BreakerConfig struct {
    Threshold      int           // 错误阈值
    FailureRate    float64       // 错误率阈值
    Timeout        time.Duration // 熔断超时时间
    HalfOpenLimit  int           // 半开状态请求限制
}

// 熔断器
type CircuitBreaker struct {
    name           string
    config         *BreakerConfig
    state          CircuitState
    failures       int
    totalRequests  int
    lastFailure    time.Time
    halfOpenCount  int
    mu             sync.RWMutex
    logger         *LoggerService
}

func NewCircuitBreaker(name string, config *BreakerConfig, logger *LoggerService) *CircuitBreaker {
    return &CircuitBreaker{
        name:    name,
        config:  config,
        state:   StateClosed,
        logger:  logger,
    }
}

// 执行受保护的函数
func (cb *CircuitBreaker) Execute(fn func() error) error {
    if !cb.allowRequest() {
        return errors.New("circuit breaker is open")
    }

    err := fn()
    cb.recordResult(err)
    return err
}

// 判断是否允许请求
func (cb *CircuitBreaker) allowRequest() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()

    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        if time.Since(cb.lastFailure) > cb.config.Timeout {
            cb.mu.RUnlock()
            cb.mu.Lock()
            cb.state = StateHalfOpen
            cb.halfOpenCount = 0
            cb.mu.Unlock()
            cb.mu.RLock()
            return true
        }
        return false
    case StateHalfOpen:
        return cb.halfOpenCount < cb.config.HalfOpenLimit
    default:
        return false
    }
}

// 记录执行结果
func (cb *CircuitBreaker) recordResult(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    cb.totalRequests++

    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()

        // 计算错误率
        failureRate := float64(cb.failures) / float64(cb.totalRequests)

        switch cb.state {
        case StateClosed:
            if cb.failures >= cb.config.Threshold || failureRate >= cb.config.FailureRate {
                cb.tripBreaker()
            }
        case StateHalfOpen:
            cb.tripBreaker()
        }
    } else {
        switch cb.state {
        case StateHalfOpen:
            cb.halfOpenCount++
            if cb.halfOpenCount >= cb.config.HalfOpenLimit {
                cb.reset()
            }
        }
    }
}

// 触发熔断
func (cb *CircuitBreaker) tripBreaker() {
    cb.state = StateOpen
    cb.lastFailure = time.Now()
    
    cb.logger.LogSystem(
        "warn",
        "circuit_breaker",
        "trip",
        "Circuit breaker tripped",
        map[string]interface{}{
            "name":         cb.name,
            "failures":     cb.failures,
            "total":        cb.totalRequests,
            "failure_rate": float64(cb.failures)/float64(cb.totalRequests),
        },
    )
}

// 重置熔断器
func (cb *CircuitBreaker) reset() {
    cb.state = StateClosed
    cb.failures = 0
    cb.totalRequests = 0
    cb.halfOpenCount = 0
    
    cb.logger.LogSystem(
        "info",
        "circuit_breaker",
        "reset",
        "Circuit breaker reset",
        map[string]interface{}{
            "name": cb.name,
        },
    )
} 