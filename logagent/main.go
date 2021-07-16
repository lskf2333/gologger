package main

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"goLoggerTest/logagent/etcd"
	"goLoggerTest/logagent/kafka"
	"goLoggerTest/logagent/tailfile"
	"gopkg.in/ini.v1"
	"strings"
	"time"
)

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcdConfig    `ini:"etcd"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int64  `ini:"chan_size"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

type EtcdConfig struct {
	Address string `ini:"address"`
}

//日志收集的客户端
//类似的开源项目还有filebeat
//手机指定目录下的日志文件，发送到kafka中

func run() (err error) {
	//循环读数据
	for {
		line, ok := <-tailfile.TailObj.Lines //chan tail.Line
		if !ok {
			logrus.Warn("tail file close reope,filename:%s\n", tailfile.TailObj.Filename)
			time.Sleep(time.Second)
			continue
		}
		//如果是空行，就跳过
		if len(strings.TrimSpace(line.Text)) == 0 {
			logrus.Info("出现空行了，跳过")
			continue
		}
		//利用通道将同步的代码改为异步的
		//把都出来的一行日志保诚成Kafka里面的msg类型，丢到通道中
		msg := &sarama.ProducerMessage{}
		msg.Topic = "web_log"
		msg.Value = sarama.StringEncoder(line.Text)
		kafka.ToMsgChan(msg)
	}
	return
}

func main() {
	var configObj = new(Config)
	// 0.读配置文件 `go-ini`

	// cfg, err := ini.Load("./conf/config.ini")
	// if err != nil {
	// 	logrus.Error("load config failed,err:%v", err)
	// 	return
	// }
	// kafkaAddr := cfg.Section("kafka").Key("address").String()
	// fmt.Println(kafkaAddr)

	err := ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Error("load config failed,err:%v", err)
		return
	}

	// 1.初始化连接kafka（做好准备工作）
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Error("init kafka failed,err:%v", err)
		return
	}
	logrus.Info("init kafka success!")

	//初始化etcd连接
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Error("init etcd failed,err:%v", err)
		return
	}
	logrus.Info("init etcd success!")
	//从etcd中拉去要手机日志的配置项

	// 2.根据配置中的日志路径初始化tail包
	err = tailfile.Init(configObj.CollectConfig.LogFilePath)
	if err != nil {
		logrus.Error("init tailfile failed,err:%v", err)
		return
	}
	logrus.Info("init tailfile success!")
	// 3.把日志用过sarama发送kafka
	err = run()
	if err != nil {
		logrus.Error("run failed,err:%v", err)
		return
	}
}
