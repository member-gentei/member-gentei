package apis

import (
	"fmt"
	"net/url"

	"github.com/otiai10/opengraph/v2"
)

func GetYouTubeChannelOG(channelID string) (*opengraph.OpenGraph, error) {
	channelURL := fmt.Sprintf("https://youtube.com/channel/%s", channelID)
	return GetOpenGraph(channelURL)
}

// Gets OpenGraph data from the provided URL.
func GetOpenGraph(uri string) (*opengraph.OpenGraph, error) {
	if _, err := url.Parse(uri); err != nil {
		return nil, err
	}
	return opengraph.Fetch(uri)
}
