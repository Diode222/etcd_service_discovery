package etcdservice
import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"strings"
	"time"

)

// Register register service with name as prefix to etcd, multi etcd addr should use ; to split
func register(etcdAddr, name string, addr string, ttl int64) error {
	var err error

	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(etcdAddr, ";"),
			DialTimeout: 15 * time.Second,
		})
		if err != nil {
			fmt.Printf("connect to etcd err:%s", err)
			return err
		}
	}

	ticker := time.NewTicker(time.Second * time.Duration(ttl))

	go func() {
		for {
			getResp, err := cli.Get(context.Background(), "/"+schema+"/"+name+"/"+addr)
			//fmt.Printf("getResp:%+v\n",getResp)
			if err != nil {
				log.Println(err)
				fmt.Printf("Register:%s", err)
			} else if getResp.Count == 0 {
				err = withAlive(name, addr, ttl)
				if err != nil {
					log.Println(err)
					fmt.Printf("keep alive:%s", err)
				}
			} else {
				//fmt.Printf("getResp:%+v, do nothing\n",getResp)
			}

			<-ticker.C
		}
	}()

	return nil
}

func withAlive(name string, addr string, ttl int64) error {
	leaseResp, err := cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	//fmt.Printf("key:%v\n", "/"+schema+"/"+name+"/"+addr)
	_, err = cli.Put(context.Background(), "/"+schema+"/"+name+"/"+addr, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		fmt.Printf("put etcd error:%s",err)
		return err
	}

	_, err = cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		fmt.Printf("keep alive error:%s",err)
		return err
	}
	return nil
}

// UnRegister remove service from etcd
func unRegister(name string, addr string) {
	if cli != nil {
		cli.Delete(context.Background(), "/"+schema+"/"+name+"/"+addr)
	}
}
