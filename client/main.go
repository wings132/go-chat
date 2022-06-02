package main

import (
	"fmt"
	"go-chat/client/logger"
	"go-chat/client/process"
	"google.golang.org/grpc"
	"log"

	gp "github.com/howeyc/gopass"
)

func main() {
	var (
		key             int
		loop            = true
		userName        string
		password        string
	)

	conn, err := grpc.Dial("127.0.0.1:8889", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can't connect server: %v", err)
		return
	}
	defer conn.Close()

	for loop {
		logger.Info("\n----------------Welcome to the chat room--------------\n")
		logger.Info("\t\tSelect the options：\n")
		logger.Info("\t\t\t 1、Sign in\n")
		logger.Info("\t\t\t 2、Sign up\n")
		logger.Info("\t\t\t 3、Exit the system\n")

		// get user input
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			logger.Info("sign In Please\r\n")
			logger.Notice("Username:")
			fmt.Scanf("%s\n", &userName)

			logger.Notice("Password:")
			// fmt.Scanf("%s\n", &password)
			ps, _ := gp.GetPasswdMasked()
			password = string(ps)

			// err := login(userName, password)
			up := process.UserProcess{Conn:conn}
			ret := up.Login(userName, password)

			if ret {
				logger.Success("Login succeed!\r\n")
			} else {
				logger.Error("Login failed\r\n")
				return
			}
		case 2:
			logger.Info("Create account\n")

			// get username
			logger.Notice("user name：")
			fmt.Scanf("%s\n", &userName)

			// get password
			logger.Notice("password：")
			//fmt.Scanf("%s\n", &password)
			ps, _ := gp.GetPasswdMasked()
			password = string(ps)

			up := process.UserProcess{Conn:conn}
			ret := up.Register(userName, password)
			if !ret {
				logger.Error("Create account failed: %v\n", err)
			}
		case 3:
			logger.Warn("Exit...\n")
			loop = false // this is equal to 'os.Exit(0)'
		default:
			logger.Error("Select is invalid!\n")
		}
	}
}
