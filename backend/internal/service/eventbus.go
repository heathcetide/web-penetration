package service

import (
	"encoding/json"
	"sync"
	"time"
)

type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time             `json:"timestamp"`
}

type EventHandler func(*Event)

type EventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
	logger   *LoggerService
}

func NewEventBus(logger *LoggerService) *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
		logger:   logger,
	}
}

// 订阅事件
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// 发布事件
func (eb *EventBus) Publish(event *Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	handlers, exists := eb.handlers[event.Type]
	if !exists {
		return
	}

	// 记录事件
	eventJSON, _ := json.Marshal(event)
	eb.logger.LogSystem(
		"info",
		"eventbus",
		"publish",
		"Event published",
		map[string]interface{}{
			"event": string(eventJSON),
		},
	)

	// 异步处理事件
	for _, handler := range handlers {
		go func(h EventHandler) {
			defer func() {
				if err := recover(); err != nil {
					eb.logger.LogSystem(
						"error",
						"eventbus",
						"handler_panic",
						"Event handler panicked",
						map[string]interface{}{
							"error": err,
							"event": event,
						},
					)
				}
			}()
			h(event)
		}(handler)
	}
} 