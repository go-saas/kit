package generate

//go:generate protoc --proto_path=. --proto_path=../../proto --proto_path=../../pkg --go_out=paths=source_relative:. internal/conf/conf.proto
