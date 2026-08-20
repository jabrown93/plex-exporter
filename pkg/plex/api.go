package plex

// Minimal Plex API types and calls previously provided by
// github.com/jrudio/go-plex-client (unmaintained). Only the fields this
// exporter reads are declared; JSON decoding ignores the rest.

import "encoding/json"

type User struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Player struct {
	Device  string `json:"device"`
	Product string `json:"product"`
}

type Part struct {
	Decision string `json:"decision"`
}

type Media struct {
	Bitrate         int    `json:"bitrate"`
	VideoResolution string `json:"videoResolution"`
	Part            []Part `json:"Part"`
}

type Metadata struct {
	Player           Player      `json:"Player"`
	User             User        `json:"User"`
	GrandparentTitle string      `json:"grandparentTitle"`
	LibrarySectionID json.Number `json:"librarySectionID"`
	ParentTitle      string      `json:"parentTitle"`
	RatingKey        string      `json:"ratingKey"`
	SessionKey       string      `json:"sessionKey"`
	Media            []Media     `json:"Media"`
	Title            string      `json:"title"`
	Type             string      `json:"type"`
}

type CurrentSessions struct {
	MediaContainer struct {
		Metadata []Metadata `json:"Metadata"`
	} `json:"MediaContainer"`
}

type MediaMetadata struct {
	MediaContainer struct {
		Metadata []Metadata `json:"Metadata"`
	} `json:"MediaContainer"`
}

type PlaySessionStateNotification struct {
	RatingKey  string `json:"ratingKey"`
	SessionKey string `json:"sessionKey"`
	State      string `json:"state"`
	ViewOffset int64  `json:"viewOffset"`
}

type NotificationContainer struct {
	PlaySessionStateNotification []PlaySessionStateNotification `json:"PlaySessionStateNotification"`
	Type                         string                         `json:"type"`
}

// GetSessions returns the currently active playback sessions.
func (c *Client) GetSessions() (CurrentSessions, error) {
	var sessions CurrentSessions
	err := c.Get("/status/sessions", &sessions)
	return sessions, err
}

// GetMetadata returns library metadata for a rating key.
func (c *Client) GetMetadata(ratingKey string) (MediaMetadata, error) {
	var metadata MediaMetadata
	err := c.Get("/library/metadata/"+ratingKey, &metadata)
	return metadata, err
}
