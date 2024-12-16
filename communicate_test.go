package edgettstool

import (
	"os"
	"testing"
)

func TestHandleGenerateTTS(t *testing.T) {
	thunk, error := NewCommunicate(DEFAULT_LANG, DEFAULT_VOICE, DEFAULT_VOLUME).HandleGenerateTTS("你好啊，见到你很高兴")

	if error != nil {
		t.Errorf("转换: %s\n", error.Error())
		return
	}

	if error := os.WriteFile("./example-1.mp3", thunk, 0777); error != nil {
		t.Errorf("写入失败: %s\n", error.Error())
		return
	}
	t.Logf("写入成功 \n")
}

func TestHandleSaveTTSFile(t *testing.T) {
	if error := NewCommunicate(DEFAULT_LANG, DEFAULT_VOICE, DEFAULT_VOLUME).HandleSaveTTSFile("你好啊，见到你很高兴", "./example-2.mp3", 0777); error != nil {
		t.Errorf("写入失败: %s\n", error.Error())
		return
	}
	t.Logf("写入成功 \n")
}
