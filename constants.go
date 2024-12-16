package edgettstool

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	BASE_URL             = "speech.platform.bing.com/consumer/speech/synthesize/readaloud"
	TRUSTED_CLIENT_TOKEN = "6A5AA1D4EAFF4E9FB37E23D68491D6F4"

	WSS_URL    = fmt.Sprintf("wss://%s/edge/v1?TrustedClientToken=%s", BASE_URL, TRUSTED_CLIENT_TOKEN)
	VOICE_LIST = fmt.Sprintf("https://%s/voices/list?trustedclienttoken=%s", BASE_URL, TRUSTED_CLIENT_TOKEN)

	DEFAULT_VOICE  = "zh-CN-XiaoyiNeural"
	DEFAULT_LANG   = "zh-CN"
	DEFAULT_VOLUME = "+40%"

	CHROMIUM_FULL_VERSION  = "130.0.2849.68"
	CHROMIUM_MAJOR_VERSION = strings.Split(CHROMIUM_FULL_VERSION, ".")[0]
	SEC_MS_GEC_VERSION     = fmt.Sprintf("1-%s", CHROMIUM_FULL_VERSION)
	WSS_HEADERS            = http.Header{
		"Host":            {"speech.platform.bing.com"},
		"User-Agent":      {fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s.0.0.0 Safari/537.36 Edg/%s.0.0.0", CHROMIUM_MAJOR_VERSION, CHROMIUM_MAJOR_VERSION)},
		"Accept-Encoding": {"gzip, deflate, br"},
		"Accept-Language": {"en-US,en;q=0.9"},
		"Pragma":          {"no-cache"},
		"Cache-Control":   {"no-cache"},
		"Origin":          {"chrome-extension://jdiccldimpdaibmpdkjnbmckianbfold"},
	}
	VOICE_HEADERS = map[string]string{
		"User-Agent":       fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s.0.0.0 Safari/537.36 Edg/%s.0.0.0", CHROMIUM_MAJOR_VERSION, CHROMIUM_MAJOR_VERSION),
		"Accept-Encoding":  "gzip, deflate, br",
		"Accept-Language":  "en-US,en;q=0.9",
		"Authority":        "speech.platform.bing.com",
		"Sec-CH-UA-Mobile": "?0",
		"Accept":           "*/*",
		"Sec-Fetch-Site":   "none",
		"Sec-Fetch-Mode":   "cors",
		"Sec-Fetch-Dest":   "empty",
		"Sec-CH-UA":        fmt.Sprintf("\" Not;A Brand\";v=\"99\", \"Microsoft Edge\";v=\"%s\", \"Chromium\";v=\"%s\"", CHROMIUM_MAJOR_VERSION, CHROMIUM_MAJOR_VERSION),
	}
)
