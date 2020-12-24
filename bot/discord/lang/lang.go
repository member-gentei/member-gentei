package lang

import (
	"encoding/base64"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// NewBundle creates the standard *i18n.Bundle. Panics.
func NewBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.AmericanEnglish)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	var err error
	_, err = bundle.ParseMessageFileBytes(mustB64Decode(enUS), "active.en-US.toml")
	if err != nil {
		panic(fmt.Errorf("error parsing active.en-US.toml: %+v", err))
	}
	_, err = bundle.ParseMessageFileBytes(mustB64Decode(zhHant), "active.zh-Hant.toml")
	if err != nil {
		panic(fmt.Errorf("error parsing active.zh-Hant.toml: %+v", err))
	}
	return bundle
}

func mustB64Decode(str string) []byte {
	dec, err := base64.StdEncoding.DecodeString(zhHant)
	if err != nil {
		panic(err)
	}
	return dec
}
