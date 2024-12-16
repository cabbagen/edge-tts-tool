package edgettstool

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Communicate struct {
	voice          string
	lang           string
	outputFormat   string
	proxy          string
	rate           string
	pitch          string
	volume         string
	receiveTimeout int
}

func NewCommunicate(lang, voice, volume string) *Communicate {
	return &Communicate{
		lang:           lang,
		voice:          voice,
		outputFormat:   "audio-24khz-48kbitrate-mono-mp3",
		proxy:          "",
		rate:           "+0%",
		pitch:          "+0Hz",
		volume:         volume,
		receiveTimeout: 10000,
	}
}

func (c *Communicate) HandleGenerateTTS(text string) ([]byte, error) {
	conn, error := c.handleWebSocketConnect()

	if error != nil {
		return nil, error
	}

	defer conn.Close()

	if error := conn.WriteMessage(websocket.TextMessage, []byte(c.getCommandRequestContent())); error != nil {
		return nil, error
	}

	if error := conn.WriteMessage(websocket.TextMessage, []byte(c.getSSMLRequestContent(text))); error != nil {
		return nil, error
	}

	return c.handleReadFromWebSocket(conn)
}

func (c *Communicate) HandleSaveTTSFile(text string, filePath string, mode os.FileMode) error {
	thunk, error := c.HandleGenerateTTS(text)

	if error != nil {
		return error
	}
	return os.WriteFile(filePath, thunk, mode)
}

func (c *Communicate) handleWebSocketConnect() (*websocket.Conn, error) {
	dialer := &websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  60 * time.Second,
		EnableCompression: true,
	}

	if len(c.proxy) > 0 {
		proxyUrl, error := url.Parse(c.proxy)

		if error != nil {
			return nil, error
		}

		dialer.Proxy = http.ProxyURL(proxyUrl)
	}

	requestUrl := fmt.Sprintf("%s&Sec-MS-GEC=%s&Sec-MS-GEC-Version=%s", WSS_URL, GenerateSecMsGec(), SEC_MS_GEC_VERSION)

	dialCtx, dialContextCancel := context.WithTimeout(context.Background(), time.Duration(c.receiveTimeout)*time.Second)

	defer func() {
		dialContextCancel()
	}()

	conn, _, error := dialer.DialContext(dialCtx, requestUrl, WSS_HEADERS)

	if error != nil {
		return nil, error
	}

	return conn, nil
}

func (c *Communicate) getCommandRequestContent() string {
	var builder strings.Builder

	builder.WriteString("Content-Type:application/json; charset=utf-8\r\n")

	builder.WriteString("Path:speech.config\r\n\r\n")

	builder.WriteString(`{"context":{"synthesis":{"audio":{"metadataoptions":{`)
	builder.WriteString(`"sentenceBoundaryEnabled":false,"wordBoundaryEnabled":true},`)
	builder.WriteString(fmt.Sprintf(`"outputFormat":"%s"`, c.outputFormat))
	builder.WriteString("}}}}\r\n")

	return builder.String()
}

func (c *Communicate) getSSMLRequestContent(text string) string {
	requestId, timestamp := strings.ReplaceAll(uuid.New().String(), "-", ""), time.Now().UTC().Format("Mon Jan 02 2006 15:04:05 GMT+0000 (Coordinated Universal Time)")

	headers := fmt.Sprintf("X-RequestId:%s\r\nContent-Type:application/ssml+xml\r\nX-Timestamp:%sZ\r\nPath:ssml\r\n\r\n", requestId, timestamp)
	ssmlTxt := fmt.Sprintf("<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='%s'><voice name='%s'><prosody pitch='%s' rate='%s' volume='%s'>%s</prosody></voice></speak>", c.lang, c.voice, c.pitch, c.rate, c.volume, text)

	return headers + ssmlTxt
}

func (c *Communicate) handleReadFromWebSocket(conn *websocket.Conn) ([]byte, error) {
	finished, failed, audioData := make(chan struct{}), make(chan error), make([]byte, 0)

	go func() {
		defer func() {
			close(failed)
			close(finished)
		}()

		for {
			receivedType, receivedData, receivedErr := conn.ReadMessage()

			// 已经断开链接
			if receivedType == -1 && receivedData == nil && receivedErr != nil {
				failed <- receivedErr
				return
			}

			switch receivedType {
			case websocket.TextMessage:
				textHeader, _, textErr := c.getHeadersAndData(receivedData)

				if textErr != nil {
					failed <- textErr
					return
				}
				if string(textHeader["Path"]) == "turn.end" {
					finished <- struct{}{}
					return
				}
			case websocket.BinaryMessage:
				if len(receivedData) < 2 {
					failed <- errors.New("we received a binary message, but it is missing the header length")
					return
				}

				headerLength := binary.BigEndian.Uint16(receivedData[:2])

				if len(receivedData) < int(headerLength+2) {
					failed <- errors.New("we received a binary message, but it is missing the audio data")
					return
				}

				audioData = append(audioData, receivedData[2+headerLength:]...)
			default:
				log.Println("recv:", receivedData)
			}
		}
	}()

	select {
	case <-finished:
		return audioData, nil
	case error := <-failed:
		return nil, error
	}
}

func (c *Communicate) getHeadersAndData(data interface{}) (map[string][]byte, []byte, error) {
	var dataBytes []byte
	switch v := data.(type) {
	case string:
		dataBytes = []byte(v)
	case []byte:
		dataBytes = v
	default:
		return nil, nil, errors.New("data must be string or []byte")
	}

	headers, headerEnd := make(map[string][]byte), bytes.Index(dataBytes, []byte("\r\n\r\n"))

	if headerEnd == -1 {
		return nil, nil, errors.New("invalid data format: no header end")
	}

	headerLines := bytes.Split(dataBytes[:headerEnd], []byte("\r\n"))

	for _, line := range headerLines {
		parts := bytes.SplitN(line, []byte(":"), 2)

		if len(parts) != 2 {
			return nil, nil, errors.New("invalid header format")
		}

		headers[string(bytes.TrimSpace(parts[0]))] = bytes.TrimSpace(parts[1])
	}

	return headers, dataBytes[headerEnd+4:], nil
}
