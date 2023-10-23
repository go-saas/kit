package conf

import "google.golang.org/protobuf/proto"

const defaultKey = "default"

func (x *PaymentConf) GetMethodOrDefault(key string) *PaymentMethod {
	var conf *PaymentMethod
	if v, ok := x.Methods[defaultKey]; ok {
		conf = v
	}
	if v, ok := x.Methods[key]; ok && v != nil {
		if conf != nil {
			proto.Merge(conf, v)
		} else {
			conf = v
		}
	}
	return conf
}
