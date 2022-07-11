package v1

import (
	"github.com/go-saas/kit/event"
)

func (x *MessageProto) ToEvent() event.Event {
	ret := event.NewMessage(x.Key, x.Value)
	for k, v := range x.Header {
		ret.Header().Set(k, v)
	}
	return ret
}

func MessageProtoFromEvent(e event.Event) *MessageProto {
	ret := &MessageProto{
		Key:    e.Key(),
		Value:  e.Value(),
		Header: map[string]string{},
	}

	for _, s := range e.Header().Keys() {
		ret.Header[s] = e.Header().Get(s)
	}
	return ret
}
