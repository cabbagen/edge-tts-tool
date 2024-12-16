package edgettstool

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

var (
	S_TO_NS   int64 = 1e9
	WIN_EPOCH int64 = 11644473600
)

func GenerateSecMsGec() string {
	ticks := (time.Now().Unix() + WIN_EPOCH)

	ticks -= ticks % 300

	ticks *= S_TO_NS / 100

	return strings.ToUpper(fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d%s", ticks, TRUSTED_CLIENT_TOKEN)))))
}
