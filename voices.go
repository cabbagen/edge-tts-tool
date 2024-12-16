package edgettstool

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type VoiceTag struct {
	ContentCategories  []string `json:"ContentCategories"`
	VoicePersonalities []string `json:"VoicePersonalities"`
}

type Voice struct {
	Name           string   `json:"Name"`
	ShortName      string   `json:"ShortName"`
	Gender         string   `json:"Gender"`
	Locale         string   `json:"Locale"`
	SuggestedCodec string   `json:"SuggestedCodec"`
	FriendlyName   string   `json:"FriendlyName"`
	Status         string   `json:"Status"`
	VoiceTag       VoiceTag `json:"VoiceTag"`
}

func GetVoiceList() ([]Voice, error) {
	response, error := resty.New().R().SetHeaders(VOICE_HEADERS).Get(VOICE_LIST)

	if error != nil {
		return nil, error
	}

	if response.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list voices, http status code: %s", response.Status())
	}

	var voices []Voice

	if error := json.Unmarshal(response.Body(), &voices); error != nil {
		return nil, error
	}

	return voices, nil
}
