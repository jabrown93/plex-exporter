package plex

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/websocket"
)

var (
	ErrAlreadyListening = errors.New("already listening")
)

type plexListener struct {
	server         *Server
	activeSessions *sessions
	log            log.Logger
}

// websocketNotification is the payload of a Plex websocket message.
type websocketNotification struct {
	NotificationContainer NotificationContainer `json:"NotificationContainer"`
}

func (s *Server) Listen(ctx context.Context, log log.Logger) error {
	s.mtx.Lock()
	if s.listener != nil {
		s.mtx.Unlock()
		return ErrAlreadyListening
	}

	s.listener = &plexListener{
		server:         s,
		activeSessions: NewSessions(ctx, s),
		log:            log,
	}

	s.mtx.Unlock()

	wsURL := *s.URL
	if wsURL.Scheme == "https" {
		wsURL.Scheme = "wss"
	} else {
		wsURL.Scheme = "ws"
	}
	wsURL.Path = "/:/websockets/notifications"

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, wsURL.String(), http.Header{
		"X-Plex-Token": []string{s.Token},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", s.URL.String(), err)
	}
	defer conn.Close()

	// Unblock the read loop when the context is cancelled.
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	level.Info(log).Log("msg", "Successfully connected", "machineID", s.ID, "server", s.Name)

	for {
		var notification websocketNotification
		if err := conn.ReadJSON(&notification); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) && closeErr.Code == websocket.CloseNormalClosure {
				return nil
			}
			level.Error(log).Log("msg", "error in websocket processing", "err", err)
			return err
		}

		if notification.NotificationContainer.Type == "playing" {
			s.listener.onPlayingHandler(notification.NotificationContainer)
		}
	}
}

func getSessionByID(sessions CurrentSessions, sessionID string) *Metadata {
	for _, session := range sessions.MediaContainer.Metadata {
		if sessionID == session.SessionKey {
			return &session
		}
	}
	return nil
}

func (l *plexListener) onPlayingHandler(c NotificationContainer) {
	err := l.onPlaying(c)
	if err != nil {
		level.Error(l.log).Log("msg", "error handling OnPlaying event", "event", c, "err", err)
	}
}

func (l *plexListener) onPlaying(c NotificationContainer) error {
	sessions, err := l.server.Client.GetSessions()
	if err != nil {
		return fmt.Errorf("error fetching sessions: %w", err)
	}

	for _, n := range c.PlaySessionStateNotification {
		if sessionState(n.State) == stateStopped {
			// When the session is stopped we can't look up the user info or media anymore.
			l.activeSessions.Update(n.SessionKey, sessionState(n.State), nil, nil)
			continue
		}

		session := getSessionByID(sessions, n.SessionKey)
		if session == nil {
			return fmt.Errorf("error getting session with key %s %+v", n.SessionKey, n)
		}

		metadata, err := l.server.Client.GetMetadata(n.RatingKey)
		if err != nil {
			return fmt.Errorf("error fetching metadata for key %s: %w", n.RatingKey, err)
		}

		level.Info(l.log).Log("msg", "Received PlaySessionStateNotification",
			"SessionKey", n.SessionKey,
			"userName", session.User.Title,
			"userID", session.User.ID,
			"state", n.State,
			"mediaTitle", metadata.MediaContainer.Metadata[0].Title,
			"mediaID", metadata.MediaContainer.Metadata[0].RatingKey,
			"timestamp", time.Duration(time.Millisecond)*time.Duration(n.ViewOffset))

		l.activeSessions.Update(n.SessionKey, sessionState(n.State), session, &metadata.MediaContainer.Metadata[0])
	}

	return nil
}
