package generate

//go:generate protoc --proto_path=. --proto_path=../../proto --go_out=paths=source_relative:. internal/conf/conf.proto
