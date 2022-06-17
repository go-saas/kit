package localize

//go:generate protoc --proto_path=../../pkg --proto_path=../../proto --go_out=paths=source_relative:../  ../../pkg/localize/i18n.proto
