package service

import (
	"encoding/json"
	"fin_im/conf"
	"fin_im/pkg/e"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (m *ClientManage) Start() {
	for {
		fmt.Println("-----监听管道通信--------")

		select {
		case conn := <-Manager.Register:
			logrus.Info("有新的链接 ", conn.ID)
			Manager.Clients[conn.ID] = conn
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: "已经链接到服务器",
			}

			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)

		case conn := <-Manager.Unregister:
			logrus.Info("连接失败", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "连接中断",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
		case broadcast := <-Manager.Broadcast:
			message := broadcast.Message
			sendId := broadcast.Client.SendId
			flag := false
			for id, conn := range Manager.Clients {
				if id != sendId {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					delete(Manager.Clients, conn.ID)
				}
			}

			id := broadcast.Client.ID

			if flag {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketOnlineRely,
					Content: "对方在线",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err := InsertMsg(conf.MongoDBName, id, string(message), 1, int64(3*month)) // 1已读
				if err != nil {
					logrus.Error("插入数据库失败")
				}
			} else {
				logrus.Info("对方不在线")
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketOfflineReply,
					Content: "对方不在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err := InsertMsg(conf.MongoDBName, id, string(message), 0, int64(3*month))
				if err != nil {
					logrus.Error(err.Error())
				}
			}
		}
	}
}
