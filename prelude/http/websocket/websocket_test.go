package websocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/hyphengolang/prelude/testing/is"
)

func testHandler(pool *Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// NOTE unsure if this should be defined here or in the logic
		if pool.IsFull() {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		rwc, err := UpgradeHTTP(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUpgradeRequired)
			return
		}

		pool.ListenAndServe(rwc)
	}
}

func TestWebSocket(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	pool := &Pool{
		Capacity: 2,
		// WriteTimeout: 10 * time.Second,
		// ReadTimeout:    10 * time.Second,
		// ReadBufferSize: 512,
	}

	srv := httptest.NewServer(http.HandlerFunc(testHandler(pool)))
	t.Cleanup(func() { srv.Close(); pool.Close() })

	t.Run("run test handler", func(t *testing.T) {

		wsURL := stripPrefix(srv.URL)
		c1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		is.NoErr(err) // connect to server

		c2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		is.NoErr(err) // connect to server

		t.Cleanup(func() { c1.Close(); c2.Close() })

		err = c1.WriteMessage(websocket.TextMessage, []byte("Hello, World!"))
		is.NoErr(err) // send message to server

		_, msg, err := c2.ReadMessage()
		is.NoErr(err) //read message from server

		is.Equal(msg, []byte("Hello, World!"))

		_, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
		is.True(err != nil) // pool is full
	})
}

var stripPrefix = func(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}
