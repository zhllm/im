package e

var codeMsg = map[Code]string{
	WebsocketSuccessMessage: "解析content内容消息",
	WebsocketSuccess:        "发送消息 请求历史记录操作成功",
	WebsocketEnd:            "请求历史记录，但是没有更多记录",
	WebsocketOfflineReply:   "针对回复消息离线应答成功",
	WebsocketOnlineRely:     "针对回复消息在线回复成功",
	WebsocketLimit:          "请求受到限制",
}
