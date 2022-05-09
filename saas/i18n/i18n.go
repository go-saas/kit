package i18n

import (
	_ "embed"
	localize "github.com/goxiaoy/go-saas-kit/pkg/localize"
)

var (
	//go:embed  en.toml
	En []byte
	//go:embed  zh.toml
	Zh    []byte
	Files = []localize.FileBundle{
		{
			Buf: En, Path: "en.toml",
		},
		{
			Buf: Zh, Path: "zh.toml",
		},
	}
)
