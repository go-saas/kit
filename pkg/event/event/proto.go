package event

import "google.golang.org/protobuf/proto"

func NewMessageFromProto(msg proto.Message) (Event, error) {
	key, v, err := KVFromProto(msg)
	if err != nil {
		return nil, err
	}
	return NewMessage(key, v), nil
}

func KVFromProto(msg proto.Message) (string, []byte, error) {
	key := string(msg.ProtoReflect().Descriptor().FullName())
	data, err := proto.Marshal(msg)
	if err != nil {
		return key, nil, err
	}
	return key, data, nil
}

func ProtoHandler[T proto.Message](msg T, next HandlerOf[T]) Handler {
	key := string(msg.ProtoReflect().Descriptor().FullName())
	return FilterKeyHandler(key, TransformHandler(func(event Event) (T, error) {
		err := proto.Unmarshal(event.Value(), msg)
		return msg, err
	}, next))
}
