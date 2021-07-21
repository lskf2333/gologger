package main

//etcd 简单使用
import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	defer cli.Close()

	//put 存值
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	str := `[{"path":"d:/logs/s4.log","topic":"s4_log"},{"path":"d:/logs/web.log","topic":"web_log"}]`
	_, err = cli.Put(ctx, "collect_log_conf", str)
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}

	//get 取值
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	gr, err := cli.Get(ctx, "collect_log_conf")
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range gr.Kvs {
		fmt.Printf("key:%s value:%s\n", ev.Key, ev.Value)
	}
	cancel()
}
