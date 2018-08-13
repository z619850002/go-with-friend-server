package Utils

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"encoding/binary"
	"github.com/gorilla/websocket"
	"go-with-friend-server/src/Server/Service/RegistAndLogin"
	"go-with-friend-server/src/Server/Service/MainPageInfo"
	"go-with-friend-server/src/Server/Service/Image"
	"go-with-friend-server/src/Server/Service/AI"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}


func analyse(m []byte, c *Client) []byte{
	//get the protocol and the message body
	temp := make([]byte, 2)
	message := make([]byte, len(m)-2)
	copy(temp,m[:2])
	copy(message,m[2:])
	top := binary.BigEndian.Uint16(temp)
	//handle the request
	var replaymessage []byte
	switch top {
	case 0:		//登陆请求
		replaymessage = RegistAndLogin.Login(message)
	case 2:		//注册请求
		replaymessage = RegistAndLogin.Regist(message)
	case 4:		//请求获得玩家信息
		replaymessage = RegistAndLogin.GetPlayerInfo(message)
	case 6:		//请求获得theirturn的信息
		replaymessage = RegistAndLogin.GetTheirturn(message)
	case 8: 	//请求获得pending
		replaymessage = RegistAndLogin.GetPending(message)
	case 10:	//请求获得myturn
		replaymessage = RegistAndLogin.GetMyturn(message)
	case 14:	//请求获得对战历史
		replaymessage = MainPageInfo.GetMyHistory(message)
	case 16:	//请求获得收藏的对战历史
		replaymessage = MainPageInfo.GetCollectedHistory(message)
	case 18:	//请求获得头像
		replaymessage = Image.DownloadImg(message)
	case 20:	//ai battle请求
		AI.Startaigame(message, c.send)
	case 21:	//ai battle对战中玩家下棋位置发送
		AI.Setplayerposition(message)


	default:
		fmt.Printf("the route is null")
	}
	return replaymessage
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	fmt.Println("in func : readPump")
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)		//sets the maximum size for a message read from the peer
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		//func TrimSpace(s []byte) []byte : to cut all the blank at the beginning and the ending of s
		//func Replace(s, old, new []byte, n int) []byte : replace all newline('/n') with space(' '), n means how many newline to replace, n<0 : no limits
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		//write the message to the client.send and return to the client
		c.send <- analyse(message, c)
	}
}



// writePump pumps messages from the hub to the websocket connection.
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	//升级为websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}