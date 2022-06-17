package conf

//go:generate protoc --proto_path=../../pkg --proto_path=../../proto --go_out=paths=source_relative:../  ../../pkg/conf/conf.proto
