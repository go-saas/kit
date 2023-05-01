package multilingual

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

func TestLanguageMake(t *testing.T) {
	tag := language.Make("zh-CN")
	assert.Equal(t, "zh-CN", tag.String())
	logTag(t, tag)

	tag = language.Make("zh-TW")
	assert.Equal(t, "zh-TW", tag.String())
	logTag(t, tag)

	tag = language.Make("zh-Hans")
	assert.Equal(t, "zh-Hans", tag.String())
	logTag(t, tag)

	tag = language.Make("zh-Hant")
	assert.Equal(t, "zh-Hant", tag.String())
	logTag(t, tag)
}

func logTag(t *testing.T, tag language.Tag) {
	for {
		t.Logf("%s", tag.String())
		if tag == language.Und {
			break
		}
		tag = tag.Parent()
	}
}

func TestGetTranslation(t *testing.T) {
	lan := &Language{
		Code: "zh-CN",
		Name: "中文",
		Trans: []*LanguageTrans{
			{
				Embed: Embed{LanguageCode: language.SimplifiedChinese.String()},
				Name:  "中文",
			},
			{
				Embed: Embed{LanguageCode: language.Japanese.String()},
				Name:  "中国語",
			},
			{
				Embed: Embed{LanguageCode: language.TraditionalChinese.String()},
				Name:  "繁体中文",
			},
			{
				Embed: Embed{LanguageCode: language.English.String()},
				Name:  "Chinese",
			},
		},
	}

	trans, ok := GetTranslation(lan, []language.Tag{language.Chinese, language.SimplifiedChinese, language.English}...)
	assert.Equal(t, true, ok)
	assert.Equal(t, "中文", trans.(*LanguageTrans).Name)

	trans, ok = GetTranslation(lan, []language.Tag{language.Japanese, language.Chinese, language.SimplifiedChinese, language.English}...)
	assert.Equal(t, true, ok)
	assert.Equal(t, "中国語", trans.(*LanguageTrans).Name)

	trans, ok = GetTranslation(lan, []language.Tag{language.English}...)
	assert.Equal(t, true, ok)
	assert.Equal(t, "Chinese", trans.(*LanguageTrans).Name)

	trans, ok = GetTranslation(lan, []language.Tag{language.AmericanEnglish}...)
	assert.Equal(t, true, ok)
	assert.Equal(t, "Chinese", trans.(*LanguageTrans).Name)

	trans, ok = GetTranslation(lan, []language.Tag{language.Bengali}...)
	assert.Equal(t, true, ok)
	assert.Equal(t, "Chinese", trans.(*LanguageTrans).Name)

}
