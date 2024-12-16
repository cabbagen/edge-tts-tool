package edgettstool

import (
	"testing"
)

func TestGetVoiceList(t *testing.T) {
	voices, error := GetVoiceList()

	if error != nil {
		t.Logf("获取声音失败: %s\n", error.Error())
		return
	}
	for _, voice := range voices {
		t.Logf("Local: %s, ShortName: %s, Gender: %s\n", voice.Locale, voice.ShortName, voice.Gender)
	}

	t.Logf("获取成功")
}
