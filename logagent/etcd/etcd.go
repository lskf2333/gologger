package etcd

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

//etcd 相关操作

type collectEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

var (
	client *clientv3.Client
)

func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logrus.Error("connect to etcd failed, err:%v\n", err)
		return
	}
	return
}

//拉取日志收集配置项的函数
func GetConf(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Error("get conf from etcd by key:%s failed,err:%v\n", key, err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warningf("get len:0 conf from etcd by key:%s \n", key)
		return
	}
	ret :=resp.Kvs[0]
	var
	json.Unmarshal(ret.Value,) //json格式字符串
}
