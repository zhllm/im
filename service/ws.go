package service

import (
	"encoding/json"
	"fin_im/cache"
	"fin_im/conf"
	"fin_im/pkg/e"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const month = 60 * 60 * 24 * 30

type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

type Client struct {
	ID     string
	SendId string
	Socket *websocket.Conn
	Send   chan []byte
}

type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

type ClientManage struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

type Message struct {
	Sender    string `json:"sender, omitempty"`
	Recipient string `json:"recipient, omitempty"`
	Content   string `json:"content, omitempty"`
}

var Manager = &ClientManage{
	Clients:    make(map[string]*Client),
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}

func CreateID(uid string, toUid string) string {
	return uid + "->" + toUid //  1 -> 2
}

func Handler(c *gin.Context) {
	uid := c.Query("uid")
	toUid := c.Query("toUid")
	logrus.Info("新的链接", uid, "  ", toUid)

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		logrus.Error(err.Error())
		http.NotFound(c.Writer, c.Request)
		return
	}

	client := &Client{
		ID:     CreateID(uid, toUid),
		SendId: CreateID(toUid, uid),
		Socket: conn,
		Send:   make(chan []byte),
	}

	Manager.Register <- client
	go client.Read()
	go client.Write()
}

func (client *Client) Read() {
	defer func() {
		Manager.Unregister <- client
		_ = client.Socket.Close()
	}()

	for {

		logrus.Info("开始循环监听")

		client.Socket.PingHandler()
		sendMsg := new(SendMsg)
		err := client.Socket.ReadJSON(&sendMsg)

		logrus.Info("获取到客户端发送的消息")
		if err != nil {
			logrus.Info("accept message error", err.Error())
			_ = client.Socket.Close()
			break
		}

		if sendMsg.Type == 1 { // 1-> 2 发消息
			r1, _ := cache.RedisClient.Get(client.ID).Result()
			r2, _ := cache.RedisClient.Get(client.SendId).Result()

			if r1 > "3" && r2 == "" { // 1 给 2 发了超过3条，但是2米有回复，避免骚扰停止发送消息
				replyMsg := ReplyMsg{
					Code:    e.WebsocketLimit,
					Content: "达到限制",
				}

				msg, _ := json.Marshal(replyMsg)
				_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			} else {
				cache.RedisClient.Incr(client.ID)
				cache.RedisClient.Expire(client.ID, time.Hour*24*30*3).Result()
			}

			Manager.Broadcast <- &Broadcast{
				Client:  client,
				Message: []byte(sendMsg.Content),
			}
		} else if sendMsg.Type == 2 {
			// 获取历史消息
			timeT, err := strconv.Atoi(sendMsg.Content)
			if err != nil {
				timeT = 99999
			}

			results, _ := FindMany(conf.MongoDBName, client.SendId, client.ID, int64(timeT), 10000)

			if len(results) > 10 {
				// results = results[:10]
			} else if len(results) == 0 {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "到底了",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: result.Msg,
				}
				msg, _ := json.Marshal(replyMsg)
				_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
			}

			replyMsg := ReplyMsg{
				From:    "system",
				Content: strconv.Itoa(len(results)),
			}
			msg, _ := json.Marshal(replyMsg)
			_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func (client *Client) Write() {
	defer func() {
		_ = client.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				_ = client.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccessMessage,
				Content: fmt.Sprintln("%s", string(message)),
			}

			msg, _ := json.Marshal(replyMsg)
			_ = client.Socket.WriteMessage(websocket.TextMessage, msg)

		}
	}
}
