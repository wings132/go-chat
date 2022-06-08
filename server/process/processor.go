package process

import (
	"context"
	"fmt"
	common "go-chat/common/message"
	pb "go-chat/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"net"
	"time"
)

type Processor struct {
	ClientConnsMap map[string]string
	ser grpc.Server
	pb.UnimplementedChatServiceServer
}

func (processor Processor) OnLoginReq(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	remoteCon, _ := peer.FromContext(ctx)
	log.Printf("username: %s, password: %s,  from %s \n", req.Username, req.Password, remoteCon.Addr.String())

	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection 测试连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.ListDatabases(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := client.Database("test").Collection("user")
	var result bson.M
	err = collection.FindOne(ctx, bson.D{{"userName", req.Username},{"passWord",req.Password}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title\n")
		return &pb.LoginRes{ Result: pb.Status_FAIL}, nil
	}
	if err != nil {
		panic(any(err))
	}


	//clientConn := model.ClientConn{}
	//rand.Seed(time.Now().UnixNano())
	//clientConn.Save(req.Username,remoteCon.Addr.String())
	//clientConn.ShowAllUsers()

	processor.ClientConnsMap[req.Username] = remoteCon.Addr.String()

	return &pb.LoginRes{ Result: pb.Status_OK}, nil
}

func (processor Processor) OnRegisterReq(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {

	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection 测试连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.ListDatabases(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := client.Database("test").Collection("user")
	result, err := collection.InsertOne(ctx, bson.D{{"userName", req.Username},{"passWord",req.Password}})
	if err != nil {
		fmt.Printf("register err %s\n",err.Error())
		return &pb.RegisterRes{Result:pb.Status_FAIL,Reason: err.Error()},nil
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

	return &pb.RegisterRes{},nil
}

func (processor Processor) OnShowAllUserOnline(context.Context, *pb.ShowAllUserOnlineReq) (*pb.ShowAllUserOnlineRes, error){
	res := pb.ShowAllUserOnlineRes{}
	for k,_ := range processor.ClientConnsMap{
		res.Users = append(res.Users, k)
	}
	processor.ser.GetServiceInfo()
	return &res, nil
}

func (processor Processor) OnP2PChatReq(req *pb.P2PChatReq, srv pb.ChatService_OnP2PChatReqServer) error{

}

// 处理消息
// 根据消息的类型，使用对应的处理方式
func (this *Processor) messageProcess(message common.Message) (err error) {
	//switch message.Type {
	//case common.LoginMessageType:
	//	up := UserProcess{Conn: this.Conn}
	//	err = up.UserLogin(message.Data)
	//	if err != nil {
	//		fmt.Printf("some error: %v\n", err)
	//	}
	//case common.RegisterMessageType:
	//	up := UserProcess{Conn: this.Conn}
	//	err = up.UserRegister(message.Data)
	//	if err != nil {
	//		fmt.Printf("some error when register: %v\n", err)
	//	}
	//case common.UserSendGroupMessageType:
	//	fmt.Println("user send group message!")
	//	gmp := GroupMessageProcess{}
	//	gmp.sendToGroupUsers(message.Data)
	//case common.ShowAllOnlineUsersType:
	//	olP := OnlineInfoProcess{this.Conn}
	//	err = olP.showAllOnlineUserList()
	//	if err != nil {
	//		fmt.Printf("get all online user list error: %v\n", err)
	//	}
	//case common.PointToPointMessageType:
	//	pop := PointToPointMessageProcess{}
	//	err = pop.sendMessageToTargetUser(this.Conn, message.Data)
	//default:
	//	fmt.Printf("other type\n")
	//}
	return nil
}

func (this *Processor) MainProcess(){
	lis, err := net.Listen("tcp", "127.0.0.1:8889")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterChatServiceServer(s, Processor{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// MainProcess 处理和用户的之间的通讯
func (this *Processor) MainProcess2() {

	// 循环读来自客户端的消息
	//for {
	//	dispatcher := utils.Dispatcher{Conn: this.Conn}
	//	message, err := dispatcher.ReadData()
	//	if err != nil {
	//		if err == io.EOF {
	//			cc := model.ClientConn{}
	//			cc.Del(this.Conn)
	//			fmt.Printf("client closed!\n")
	//			break
	//		}
	//		fmt.Printf("get login message error: %v", err)
	//	}
	//
	//	fmt.Println("sender is ", this.Conn.RemoteAddr())
	//
	//	// 处理来客户端的消息
	//	// 按照消息的类型，使用不同的处理方法
	//	err = this.messageProcess(message)
	//	if err != nil {
	//		fmt.Printf("some error: %v\n", err)
	//		break
	//	}
	//}
}
