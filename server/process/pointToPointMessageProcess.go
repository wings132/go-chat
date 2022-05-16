package process

import (
	"encoding/json"
	"fmt"
	common "go-chat/common/message"
	"go-chat/server/model"
	"go-chat/server/utils"
	"net"
)

type PointToPointMessageProcess struct{}

func (this PointToPointMessageProcess) sendMessageToTargetUser(sender net.Conn, message string) (err error) {
	var pointToPointMessage common.PointToPointMessage
	err = json.Unmarshal([]byte(message), &pointToPointMessage)
	if err != nil {
		return
	}

	clientConn := model.ClientConn{}
	conn, err := clientConn.SearchByUserName(pointToPointMessage.TargetUserName)
	if err != nil {
		return
	}

	fmt.Println("conn.RemoteAddr() = ", conn.RemoteAddr())
	fmt.Println("sender.RemoteAddr() = ", sender.RemoteAddr())
	if conn.RemoteAddr() == sender.RemoteAddr() {
		fmt.Println("send msg to someone self")
	}

	var responseMessage common.ResponseMessage
	responseMessage.Type = common.PointToPointMessageType

	var responseMessageData = common.PointToPointMessage{
		SourceUserName: pointToPointMessage.SourceUserName,
		TargetUserName: pointToPointMessage.TargetUserName,
		Content:        pointToPointMessage.Content,
	}

	data, err := json.Marshal(responseMessageData)
	if err != nil {
		return
	}
	responseMessage.Data = string(data)

	responseMessage.Code = 200

	fmt.Println("send data=",responseMessage)
	responseData, err := json.Marshal(responseMessage)
	if err != nil {
		return
	}

	dispatcher := utils.Dispatcher{Conn: conn}
	err = dispatcher.WriteData(responseData)

	return
}

func (this *PointToPointMessageProcess) responseClient(conn net.Conn, code int, data string, popErr error) (err error) {
	responseMessage := common.ResponseMessage{
		Code:  code,
		Type:  common.PointToPointMessageType,
		Error: fmt.Sprintf("%v", popErr),
		Data:  data,
	}

	fmt.Println("responseClient data = ",responseMessage)
	responseData, err := json.Marshal(responseMessage)
	if err != nil {
		fmt.Printf("some error when generate response message, error: %v", err)
	}

	dispatcher := utils.Dispatcher{Conn: conn}
	err = dispatcher.WriteData(responseData)

	return
}
