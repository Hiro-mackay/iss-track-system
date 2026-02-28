package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wsReadLimit  = 512
	wsPongWait   = 60 * time.Second
	wsPingPeriod = (wsPongWait * 9) / 10
)

// HandleWebSocket returns an http.HandlerFunc that streams ISS position via WebSocket.
func HandleWebSocket(cache *TLECache, interval time.Duration, allowedOrigins []string) http.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = struct{}{}
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true
			}
			_, ok := allowed[origin]
			return ok
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("websocket upgrade failed", "error", err)
			return
		}
		defer conn.Close()

		conn.SetReadLimit(wsReadLimit)
		_ = conn.SetReadDeadline(time.Now().Add(wsPongWait))
		conn.SetPongHandler(func(string) error {
			return conn.SetReadDeadline(time.Now().Add(wsPongWait))
		})

		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		pingTicker := time.NewTicker(wsPingPeriod)
		defer pingTicker.Stop()

		sendPosition := func() bool {
			name, l1, l2, err := cache.Get()
			if err != nil {
				slog.Error("websocket: failed to get TLE", "error", err)
				return true
			}

			pos, err := GetCurrentPosition(name, l1, l2)
			if err != nil {
				slog.Error("websocket: failed to compute position", "error", err)
				return true
			}

			if err := conn.WriteJSON(pos); err != nil {
				slog.Info("websocket client disconnected", "error", err)
				return false
			}
			return true
		}

		// Send initial position immediately.
		if !sendPosition() {
			return
		}

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if !sendPosition() {
					return
				}
			case <-pingTicker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}
