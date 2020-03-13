package websocket

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
)

var (
	keyGUID                = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	ErrMethodNotAllowed    = errors.New("bat method")
	ErrBadWebsocketVersion = errors.New("missing or bad websocket version")
	ErrNotWebsocket        = errors.New("not websocket protocol")
	ErrNotHijacker         = errors.New("http.ResponseWriter does not implement http.Hijacker")
	// ErrChallengeResponse   = errors.New("mismatch challenge/response")
	ErrHttpProtoAtLeast = errors.New("http proto at least 1.1")
	ErrWebsocketKey     = errors.New("WebSocket protocol violation: missing Sec-WebSocket-Key")
)

func Upgrade(writer http.ResponseWriter, req *http.Request) (*Conn, error) {

	if errCode, err := verifyRequest(req); err != nil {
		http.Error(writer, err.Error(), errCode)
		return nil, err
	}
	key := req.Header.Get("Sec-WebSocket-Key")
	writer.Header().Set("Upgrade", "websocket")
	writer.Header().Set("Connection", "Upgrade")
	writer.Header().Set("Sec-WebSocket-Accept", secWebsocketKey(key))
	hj, ok := writer.(http.Hijacker)
	if !ok {
		http.Error(writer, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return nil, ErrNotHijacker
	}
	netConn, brw, err := hj.Hijack()
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil, err
	}
	b, _ := brw.Reader.Peek(brw.Reader.Buffered())
	brw.Reader.Reset(io.MultiReader(bytes.NewReader(b), netConn))
	return newConn(netConn, brw.Reader, brw.Writer), nil
}

func verifyRequest(req *http.Request) (int, error) {
	if !req.ProtoAtLeast(1, 1) {
		return http.StatusUpgradeRequired, ErrHttpProtoAtLeast
	}
	if req.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, ErrMethodNotAllowed
	}
	if req.Header.Get("Sec-Websocket-Version") != "13" {
		return http.StatusBadRequest, ErrBadWebsocketVersion
	}
	if req.Header.Get("Upgrade") != "websocket" {
		return http.StatusUpgradeRequired, ErrNotWebsocket
	}
	if req.Header.Get("Connection") != "Upgrade" {
		return http.StatusUpgradeRequired, ErrNotWebsocket
	}
	if req.Header.Get("Sec-WebSocket-Key") == "" {
		return http.StatusBadGateway, ErrWebsocketKey
	}
	return 0, nil
}

func secWebsocketKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	h.Write(keyGUID)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
