package websocket

import (
	"math"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var defaultUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func UpgradeHTTP(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return defaultUpgrader.Upgrade(w, r, nil)
}

type Pool struct {
	// maximum cap clients allowed
	Capacity uint
	// Maximum message size allowed from peer.
	ReadBufferSize int64
	// Time allowed to read the next pong message from the peer.
	ReadTimeout time.Duration
	// Time allowed to write a message to the peer.
	WriteTimeout time.Duration
	// container for connected channels
	cs map[*channel]struct{}

	send chan Packet
	r, d chan *channel

	isListening bool
}

var defaultPool = &Pool{
	Capacity: math.MaxInt,
}

func ListenAndServe(rwc *websocket.Conn) {
	defaultPool.ListenAndServe(rwc)
}

func (pool *Pool) IsFull() bool {
	if pool.Capacity == 0 {
		return false
	}

	return len(pool.cs) >= int(pool.Capacity)
}

func (pool *Pool) Listen() {
	pool.cs = make(map[*channel]struct{})
	pool.r = make(chan *channel)
	pool.d = make(chan *channel)
	pool.send = make(chan Packet)
	pool.isListening = true

	go listen(pool)
}

func (pool *Pool) Serve(rwc *websocket.Conn) {
	ch := &channel{send: make(chan Packet), pool: pool, rwc: rwc}
	pool.r <- ch

	go read(ch)
	go write(ch)
}

func (pool *Pool) ListenAndServe(rwc *websocket.Conn) {
	if !pool.isListening {
		pool.Listen()
	}

	pool.Serve(rwc)
}

func listen(pool *Pool) {
	for {
		select {
		case ch := <-pool.r:
			// NOTE decision to control capacity on the application level
			pool.cs[ch] = struct{}{}
		case ch := <-pool.d:
			if _, ok := pool.cs[ch]; ok {
				delete(pool.cs, ch)
				close(ch.send)
			}
		case p := <-pool.send:
			for ch := range pool.cs {
				select {
				case ch.send <- p:
				default:
					close(ch.send)
					delete(pool.cs, ch)
				}
			}
		}
	}
}

func (pool *Pool) Close() {
	close(pool.send)
	close(pool.r)
	close(pool.d)

	pool.isListening = false
}

type channel struct {
	send chan Packet
	pool *Pool

	rwc *websocket.Conn
}

func (ch *channel) respond(typ int, data []byte) error {
	return ch.rwc.WriteMessage(typ, data)
}

func (ch *channel) respondErr(code int, reason error) error {
	var msg []byte
	if reason != nil {
		msg = websocket.FormatCloseMessage(code, reason.Error())
	} else {
		msg = websocket.FormatCloseMessage(code, "")
	}
	return ch.respond(websocket.CloseMessage, msg)
}

func read(ch *channel) {
	defer func() {
		ch.pool.d <- ch
		ch.rwc.Close()
	}()

	// If read buffer size is zero then this is ignored
	if ch.pool.ReadBufferSize != 0 {
		ch.rwc.SetReadLimit(ch.pool.ReadBufferSize)
	}

	// If read timeout is zero then this is ignored
	if ch.pool.ReadTimeout != 0 {
		ch.rwc.SetReadDeadline(time.Now().Add(ch.pool.ReadTimeout))
	}

	for {
		typ, p, err := ch.rwc.ReadMessage()
		if err != nil {
			ch.respondErr(websocket.CloseInternalServerErr, nil)
			break
		}

		ch.pool.send <- Packet{typ, p}
	}
}

func write(ch *channel) {
	defer func() {
		// ticker.Stop()
		ch.rwc.Close()
	}()

	for {
		select {
		case p, ok := <-ch.send:
			// If read timeout is zero then this is ignored
			// NOTE no cleaner way of declaring this
			if ch.pool.WriteTimeout != 0 {
				ch.rwc.SetWriteDeadline(time.Now().Add(ch.pool.WriteTimeout))
			}

			if !ok {
				// the Pool closed ch.send
				ch.respondErr(websocket.CloseInternalServerErr, nil)
				return
			}

			typ, data := p.parse()
			if err := ch.rwc.WriteMessage(typ, data); err != nil {
				ch.respondErr(websocket.CloseInternalServerErr, err)
				return
			}
		}
	}
}

type Packet struct {
	typ  int
	data []byte
}

func (p Packet) parse() (int, []byte) { return p.typ, p.data }

/*
Close codes defined in RFC 6455, section 11.7.
CloseNormalClosure           = 1000
CloseGoingAway               = 1001
CloseProtocolError           = 1002
CloseUnsupportedData         = 1003
CloseNoStatusReceived        = 1005
CloseAbnormalClosure         = 1006
CloseInvalidFramePayloadData = 1007
ClosePolicyViolation         = 1008
CloseMessageTooBig           = 1009
CloseMandatoryExtension      = 1010
CloseInternalServerErr       = 1011
CloseServiceRestart          = 1012
CloseTryAgainLater           = 1013
CloseTLSHandshake            = 1015

The message types are defined in RFC 6455, section 11.8.
TextMessage denotes a text data message. The text message payload is
interpreted as UTF-8 encoded text data.
TextMessage = 1

BinaryMessage denotes a binary data message.
BinaryMessage = 2

CloseMessage denotes a close control message. The optional message
payload contains a numeric code and text. Use the FormatCloseMessage
function to format a close message payload.
CloseMessage = 8

PingMessage denotes a ping control message. The optional message payload
is UTF-8 encoded text.
PingMessage = 9

PongMessage denotes a pong control message. The optional message payload
is UTF-8 encoded text.
PongMessage = 10
*/
