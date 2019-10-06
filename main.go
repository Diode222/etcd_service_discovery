package main

import (
	"fmt"
	"github.com/Diode222/etcd_service_discovery/etcdservice"
	pb "github.com/Diode222/etcd_service_discovery/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strconv"
	"time"
)

func main() {
	serviceManager := etcdservice.NewServiceManager("127.0.0.1:2379")

	s := grpc.NewServer()
	defer s.GracefulStop()

	pb.RegisterGreeterServer(s, &server{})
	serviceManager.Register("hello_service", "127.0.0.1", "127.0.0.1", 42222, s, 5)
}

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Printf(fmt.Sprintf("%v: Receive is %s\n", time.Now(), in.GetName()))
	msg := "Hello " + in.GetName()
	return &pb.HelloReply{Message: &msg}, nil
}

func client_main() {
	serviceManger := etcdservice.NewServiceManager("127.0.0.1:2379")

	client := serviceManger.GetClient("hello_service", newGreeterClient).(pb.GreeterClient)
	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		namestr := "world " + strconv.Itoa(t.Second())
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: &namestr})
		if err == nil {
			fmt.Printf(fmt.Sprintf("%v: Reply is %s\n", t, resp.GetMessage()))
		} else {
			fmt.Printf(fmt.Sprintf("call server error:%s\n", err.Error()))
		}
	}
}

func newGreeterClient(cc *grpc.ClientConn) interface{} {
	return pb.NewGreeterClient(cc)
}
