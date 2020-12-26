package lang

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func TestNewBundle(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("caught panic creating bundle: %+v", r)
		}
	}()
	NewBundle()
}

func TestMultiLingual(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("caught panic creating bundle: %+v", r)
		}
	}()
	tlMessage := &i18n.Message{
		ID:    "BotRestartedReply",
		Other: "This bot has secretly, recently restarted and is still loading - please try again in a minute!",
	}
	bundle := NewBundle()
	enUS := i18n.NewLocalizer(bundle, "en-US")
	msg, err := enUS.LocalizeMessage(tlMessage)
	if err != nil {
		t.Logf("error translating to en-US: %+v", err)
		t.FailNow()
	}
	if msg != "This bot has secretly, recently restarted and is still loading - please try again in a minute!" {
		t.Errorf("unexpected en-US string: %s", msg)
	}
	// legitimate different language
	zhHant := i18n.NewLocalizer(bundle, "zh-Hant", "en-US")
	msg, err = zhHant.LocalizeMessage(tlMessage)
	if err != nil {
		t.Logf("error translating to zh-Hant: %+v", err)
		t.FailNow()
	}
	if msg != "這個bot最近悄悄地重啟了並且還在載入中 - 請一分鐘後再試！" {
		t.Errorf("unexpected zh-Hant string: %s", msg)
	}
	// make up a BCP47 code en-US as the fallback language
	ayyLmao := i18n.NewLocalizer(bundle, "ay-LMAO", "en-US")
	msg, err = ayyLmao.LocalizeMessage(tlMessage)
	if err != nil {
		t.Logf("error translating to ayy-LMAO: %+v", err)
		t.FailNow()
	}
	if msg != "This bot has secretly, recently restarted and is still loading - please try again in a minute!" {
		t.Errorf("unexpected ayy-LMAO (fallback to en-US) string: %s", msg)
	}
}
