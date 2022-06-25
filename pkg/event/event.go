package event

import (
	"context"
	"net/http"
)

var (
	_ Event = (*Message)(nil)
)

type Header interface {
	Get(key string) string
	Set(key string, value string)
	Keys() []string
}

type headerCarrier http.Header

// Get returns the value associated with the passed key.
func (hc headerCarrier) Get(key string) string {
	return http.Header(hc).Get(key)
}

// Set stores the key-value pair.
func (hc headerCarrier) Set(key string, value string) {
	http.Header(hc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

type Event interface {
	Header() Header
	Key() string
	Value() []byte
}

type Message struct {
	header headerCarrier
	key    string
	value  []byte
}

func (m *Message) Key() string {
	return m.key
}

func (m *Message) Header() Header {
	return m.header
}

func (m *Message) Value() []byte {
	return m.value
}

func NewMessage(key string, value []byte) Event {
	return &Message{
		key:    key,
		value:  value,
		header: headerCarrier{},
	}
}

type HandlerOf[T any] interface {
	Process(context.Context, T) error
}

type HandlerFuncOf[T any] func(context.Context, T) error

func (h HandlerFuncOf[T]) Process(ctx context.Context, e T) error {
	return h(ctx, e)
}
