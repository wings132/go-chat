package process

import (
	"encoding/json"
	"go-chat/client/utils"
	common "go-chat/common/message"
	"go-chat/config"
	"net"
)

type MessageProcess struct {
	tcpConn net.Conn
	errChan chan error
}

func (msgp *MessageProcess) getTcpConn() (conn net.Conn) {
	if msgp.tcpConn == nil {
		serverInfo := config.Configuration.ServerInfo
		msgp.tcpConn, _ = net.Dial("tcp", serverInfo.Host)
	}
	return msgp.tcpConn
}
func (msgp *MessageProcess) getErrChan() (ch chan error) {
	if msgp.errChan == nil {
		msgp.errChan = make(chan error)
	}
	return msgp.errChan
}

// user send message to server
func (msgProc MessageProcess) SendGroupMessageToServer(groupID int, userName string, content string) (err error) {
	// connect server
	serverInfo := config.Configuration.ServerInfo
	conn, err := net.Dial("tcp", serverInfo.Host)

	if err != nil {
		return
	}

	var message common.Message
	message.Type = common.UserSendGroupMessageType

	// group message
	userSendGroupMessage := common.UserSendGroupMessage{
		GroupID:  groupID,
		UserName: userName,
		Content:  content,
	}
	data, err := json.Marshal(userSendGroupMessage)
	if err != nil {
		return
	}

	message.Data = string(data)
	data, _ = json.Marshal(message)

	dispatcher := utils.Dispatcher{Conn: conn}
	err = dispatcher.SendData(data)

	return
}

// request all online user
func (msgp *MessageProcess) GetOnlineUerList() (err error) {

	var message = common.Message{}
	message.Type = common.ShowAllOnlineUsersType

	requestBody, err := json.Marshal("")
	if err != nil {
		return
	}
	message.Data = string(requestBody)

	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	dispatcher := utils.Dispatcher{Conn: msgp.getTcpConn()}
	err = dispatcher.SendData(data)
	if err != nil {
		return
	}

	// 这段代码里的channel会阻塞，暂时不知道为什么，需要深入学习channel
	// errMsg := make(chan error)
	// go Response(msg.tcpConn, errMsg)
	// err = <-errMsg

	//go Response(msgp.getTcpConn(), msgp.getErrChan())
	err = <-msgp.getErrChan()

	if err != nil {
		return
	}

	return nil
}

func (msgp *MessageProcess) PointToPointCommunication(targetUserName, sourceUserName, message string) (err error) {

	var pointToPointMessage common.Message

	pointToPointMessage.Type = common.PointToPointMessageType

	messageBody := common.PointToPointMessage{
		SourceUserName: sourceUserName,
		TargetUserName: targetUserName,
		Content:        message,
	}

	data, err := json.Marshal(messageBody)
	if err != nil {
		return
	}

	pointToPointMessage.Data = string(data)

	data, err = json.Marshal(pointToPointMessage)
	if err != nil {
		return
	}

	dispatcher := utils.Dispatcher{Conn: msgp.getTcpConn()}
	err = dispatcher.SendData(data)
	return
}
