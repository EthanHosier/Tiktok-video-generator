package captions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func msToASSTime(ms int) string {
	h := ms / 3600000
	m := (ms % 3600000) / 60000
	s := (ms % 60000) / 1000
	cs := (ms % 1000) / 10
	return fmt.Sprintf("%d:%02d:%02d.%02d", h, m, s, cs)
}

func captionsFromUrl(url string) (*Captions, error) {
	// Fetch the captions file from the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	// Read the content into memory
	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var captions Captions
	err = json.Unmarshal(byteValue, &captions)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return &captions, nil
}