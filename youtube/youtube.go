package youtube

import (
	"fmt"
	"net/url"
	"time"

	"github.com/kkdai/youtube/v2"
)

type YoutubeHandler interface {
	CaptionsForId(id string)
}

type YoutubeClient struct {
}

type Video struct {
	ID              string
	Title           string
	Description     string
	Author          string
	ChannelID       string
	ChannelHandle   string
	Views           int
	Duration        time.Duration
	PublishDate     time.Time
	DASHManifestURL string
	HLSManifestURL  string
	CaptionTrackURL string
	AudioURL        string
	VideoURL        string
}

func NewYoutubeClient() *YoutubeClient {
	return &YoutubeClient{}
}

func (y *YoutubeClient) VideoForId(id string) (*Video, error) {
	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)
	client := youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		return nil, fmt.Errorf("error getting video: %v\n\n", err)
	}

	hdFormat := video.Formats[0]
	audioFormat := video.Formats.Itag(140)[0]

	json3Url, err := replaceWithJson3Url(video.CaptionTracks[0].BaseURL)

	return &Video{
		ID:              video.ID,
		Title:           video.Title,
		Description:     video.Description,
		Author:          video.Author,
		ChannelID:       video.ChannelID,
		ChannelHandle:   video.ChannelHandle,
		Views:           video.Views,
		Duration:        video.Duration,
		PublishDate:     video.PublishDate,
		DASHManifestURL: video.DASHManifestURL,
		HLSManifestURL:  video.HLSManifestURL,
		CaptionTrackURL: json3Url,
		AudioURL:        audioFormat.URL,
		VideoURL:        hdFormat.URL,
	}, nil
}

func replaceWithJson3Url(originalURL string) (string, error) {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %v", err)
	}

	query := parsedURL.Query()
	query.Set("fmt", "json3")

	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
}
