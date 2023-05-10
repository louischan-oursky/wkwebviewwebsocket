package main

import (
	"context"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

const (
	WebsocketReadTimeout  = 30 * time.Second
	WebsocketWriteTimeout = 10 * time.Second
)

var indexHTML = `
<!DOCTYPE html>
<html>
	<head>
		<title>wkwebviewwebsocket</title>
	</head>
	<body>
		<script>
function setup() {
	let url = new URL(window.origin);
	url.protocol = "ws:";
	url.pathname = "/ws";
	let ws = new WebSocket(url);
	ws.addEventListener("open", (e) => {
		console.log("open", e);
	});
	ws.addEventListener("close", (e) => {
		console.log("close", e);
	});
	ws.addEventListener("error", (e) => {
		console.error("error", e);
	});
	ws.addEventListener("message", (e) => {
		console.log("message", e);
	});
	return ws;
}

let ws;
document.addEventListener("visibilitychange", (e) => {
	if (document.visibilityState === "visible") {
		if (ws == null || ws != null && ws.readyState === 3) {
			ws = setup();
		}
	}
	//if (document.visibilityState === "hidden") {
	//	if (ws != null && ws.readyState === 1) {
	//		ws.close();
	//		ws = null;
	//	}
	//}
});

ws = setup();
		</script>
	</body>
</html>
`

func main() {
	msgChan := make(chan []byte, 2)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(indexHTML))
	})

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		msgChan <- []byte("hello")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		rootCtx, cancel := context.WithCancel(r.Context())
		doneChan := make(chan struct{}, 2)
		errChan := make(chan error, 2)

		defer func() {
			doneChan <- struct{}{}
			doneChan <- struct{}{}
			cancel()
		}()

		wsConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
			CompressionMode:    websocket.CompressionDisabled,
		})
		if err != nil {
			return
		}
		defer wsConn.Close(websocket.StatusInternalError, "connection closed")

		go func() {
			for {
				select {
				case <-doneChan:
					return
				case n, ok := <-msgChan:
					if !ok {
						return
					}

					writeCtx, cancel := context.WithTimeout(rootCtx, WebsocketWriteTimeout)
					defer cancel()
					err := wsConn.Write(writeCtx, websocket.MessageText, n)
					if err != nil {
						errChan <- err
						return
					}
				}
			}
		}()

		go func() {
			for {
				select {
				case <-doneChan:
					return
				default:
					readCtx, cancel := context.WithTimeout(rootCtx, WebsocketReadTimeout)
					defer cancel()

					// Read everything from the connection and discard them.
					_, _, err := wsConn.Read(readCtx)
					if err != nil {
						errChan <- err
						return
					}
				}
			}
		}()

		err = <-errChan
	})

	http.ListenAndServe(":4000", nil)
}
