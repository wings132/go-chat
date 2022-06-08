package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-chat/config"
	pb "go-chat/proto"
	"go-chat/server/model"
	"go-chat/server/process"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func init() {
	// 初始化 redis 连接池，全局唯一
	redisInfo := config.Configuration.RedisInfo
	fmt.Println("redisInfo", redisInfo)
	initRedisPool(redisInfo.MaxIdle, redisInfo.MaxActive, time.Second*(redisInfo.IdleTimeout), redisInfo.Host)

	// 创建 userDao 用于操作用户信息
	// 全局唯一 UserDao 实例：model.CurrentUserDao
	model.CurrentUserDao = model.InitUserDao(pool)
}

// 和客户端的通信交互
// conn 就是客户端和服务器之间建立的连接
// 每当有个用户登陆进来之后，就启动一个 go routine
// 这个 go routine 专门用来处理服务器和客户端的通信
func dialogue(conn net.Conn) {
	//defer conn.Close()
	//processor := process.Processor{Conn: conn}
	//processor.MainProcess()
}

func dbOps() {
	// Rest of the code will go here
	// Set client options 设置连接参数
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	// Connect to MongoDB 连接数据库
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection 测试连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	databases, err := client.ListDatabases(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases.TotalSize / 1024 / 1024 / 1024)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := client.Database("test").Collection("user_sign")

	//cur, err := collection.Find(ctx, bson.D{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer cur.Close(ctx)
	//for cur.Next(ctx) {
	//	var result bson.D
	//	err := cur.Decode(&result)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("result: %v\n", result)
	//	fmt.Printf("result.Map(): %v\n\n", result.Map()["activity_id"])
	//	time.Sleep(time.Second)
	//}
	//if err := cur.Err(); err != nil {
	//	log.Fatal(err)
	//}

	// create a value into which the result can be decoded
	var result bson.M
	err = collection.FindOne(ctx, bson.D{{"uid", 160660},{"activity_id",10001}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title\n")
		return
	}
	if err != nil {
		panic(any(err))
	}


	jsondata, err := json.MarshalIndent(result,"","    ")
	if err != nil {
		panic(any(err))
	}
	fmt.Printf("%s\n", jsondata)

	fmt.Printf("Found a single document: %+v\n", result)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func main() {

	fmt.Printf("Server is already\n")

	lis, err := net.Listen("tcp", "127.0.0.1:8889")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()


	p := process.Processor{}
	p.ClientConnsMap = make(map[string]string)

	pb.RegisterChatServiceServer(s, &p)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//serverInfo := config.Configuration.ServerInfo
	//fmt.Println("serverInfo", serverInfo)
	//listener, err := net.Listen("tcp", serverInfo.Host)
	//defer listener.Close()
	//if err != nil {
	//	fmt.Printf("some error when run server, error: %v", err)
	//}
	//
	//for {
	//	fmt.Printf("Waiting for client...\n")
	//
	//	conn, err := listener.Accept()
	//	if err != nil {
	//		fmt.Printf("some error when accept server, error: %v", err)
	//	}
	//
	//	// 一旦链接成功，在启动一个协程和客户端保持通讯
	//	go dialogue(conn)
	//}
}
