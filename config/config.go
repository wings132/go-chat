package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type configuration struct {
	ServerInfo serverInfo
	RedisInfo  redisInfo
}

type serverInfo struct {
	Host string
}

type redisInfo struct {
	Host        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var Configuration = configuration{}

func init() {
	//filePath := path.Join(os.Getenv("GOPATH"), "src/go-chat/config/config.json")
	filePath := "/Users/zhaoxiang/studyFiles/openSourceProjects/go-chat/config/config.json"
	file, err := os.Open(filePath)
	if file == nil {
		fmt.Printf("open file failed")
		return
	}
	defer file.Close()
	if err != nil {
		fmt.Printf("Open file error: %v\n", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Configuration)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
