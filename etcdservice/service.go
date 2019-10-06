package etcdservice

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type ServiceManager struct {
	etcdAddr string
	resolveBuilder resolver.Builder
}

func NewServiceManager(etcdAddr string) *ServiceManager {
	builder := newResolver(etcdAddr)
	resolver.Register(builder)
	s := &ServiceManager{
		etcdAddr: etcdAddr,
		resolveBuilder: builder,
	}
	return s
}

func (s *ServiceManager) Register(serviceName string, listenIp string, serviceIp string, port int, srv *grpc.Server, ttl int64) error {
	listenAddr := fmt.Sprintf("%s:%d", listenIp, port)
	serviceAddr := fmt.Sprintf("%s:%d", serviceIp, port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln("Failed to listen, err: %s", err)
	} else {
		log.Println("Listen at: %s", serviceAddr)
	}
	defer listener.Close()

	err = register(s.etcdAddr, serviceName, serviceAddr, ttl)
	if err != nil {
		log.Println("etcd register failed.")
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		unRegister(serviceName, serviceAddr)

		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}

	}()

	if err = srv.Serve(listener); err != nil {
		log.Println("Failed to serve, err: %s", err.Error())
		panic(err)
	}

	return nil
}

func (s *ServiceManager) GetClient(serviceName string, newClientFunc func(*grpc.ClientConn) interface{}) interface{} {
	conn, err := grpc.Dial(s.resolveBuilder.Scheme()+"://author/"+ serviceName, grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		log.Println(fmt.Sprintf("GetClient failed, servicName: %s, err: %s", serviceName, err.Error()))
		panic(err)
	}

	return newClientFunc(conn)
}
