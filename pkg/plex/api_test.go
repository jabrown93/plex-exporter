package plex

import (
	"encoding/json"
	"testing"
)

func TestSessionDecodeAndLookup(t *testing.T) {
	payload := `{"MediaContainer": {"Metadata": [{
		"sessionKey": "42", "ratingKey": "1234", "type": "episode",
		"title": "Ep", "parentTitle": "Season 1", "grandparentTitle": "Show",
		"librarySectionID": 3,
		"User": {"id": "7", "title": "jared"},
		"Player": {"device": "OSX", "product": "Plex Web"},
		"Media": [{"bitrate": 4000, "videoResolution": "1080", "Part": [{"decision": "directplay"}]}]
	}]}}`

	var sessions CurrentSessions
	if err := json.Unmarshal([]byte(payload), &sessions); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if getSessionByID(sessions, "nope") != nil {
		t.Error("expected nil for unknown session key")
	}
	s := getSessionByID(sessions, "42")
	if s == nil {
		t.Fatal("expected session for key 42")
	}
	if s.User.Title != "jared" || s.Media[0].Bitrate != 4000 || s.Media[0].Part[0].Decision != "directplay" {
		t.Errorf("unexpected decode: %+v", s)
	}
	if s.LibrarySectionID.String() != "3" {
		t.Errorf("librarySectionID = %q, want 3", s.LibrarySectionID.String())
	}

	title, season, episode := labels(*s)
	if title != "Show" || season != "Season 1" || episode != "Ep" {
		t.Errorf("labels = %q %q %q", title, season, episode)
	}

	var notification websocketNotification
	msg := `{"NotificationContainer": {"type": "playing",
		"PlaySessionStateNotification": [{"sessionKey": "42", "ratingKey": "1234", "state": "paused", "viewOffset": 5000}]}}`
	if err := json.Unmarshal([]byte(msg), &notification); err != nil {
		t.Fatalf("unmarshal notification: %v", err)
	}
	n := notification.NotificationContainer
	if n.Type != "playing" || n.PlaySessionStateNotification[0].State != "paused" {
		t.Errorf("unexpected notification decode: %+v", n)
	}
}
