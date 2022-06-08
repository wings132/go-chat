package model

import (
	"fmt"
	"net"
)

type ClientConn struct{}

type ConnInfo struct {
	Conn     net.Conn
	UserName string
}

var ClientConnsMap map[string]string

func init() {
	ClientConnsMap = make(map[string]string)
}

func (cc ClientConn) Save(name string, addr string) {
	ClientConnsMap[name] = addr
}

func (cc ClientConn) Del(userConn net.Conn) {
	//for id, connInfo := range ClientConnsMap {
	//	if userConn == connInfo.Conn {
	//		delete(ClientConnsMap, id)
	//	}
	//}
}

func (cc ClientConn) SearchByUserName(userName string) (connInfo net.Conn, err error) {
	//user, err := CurrentUserDao.GetUserByUserName(userName)
	//if err != nil {
	//	return
	//}
	//
	//connInfo = ClientConnsMap[user.ID].Conn
	return
}

func (cc ClientConn)ShowAllUsers()  {
	for k, v := range ClientConnsMap{
		fmt.Printf("username[%s] addr[%s]\n",k, v)
	}
}
