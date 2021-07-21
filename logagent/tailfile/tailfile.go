package tailfile

import (
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"goLoggerTest/logagent/common"
	"goLoggerTest/logagent/kafka"
	"strings"
	"time"
)

type tailTask struct {
	path  string
	topic string
	tObj  *tail.Tail
}

func newTailTask(path, topic string) *tailTask {
	tt := &tailTask{
		path:  path,
		topic: topic,
	}
	return tt
}

func (t *tailTask) Init() (err error) {
	cfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.tObj, err = tail.TailFile(t.path, cfg)
	return
}

func (t *tailTask) run() {
	//读取日志，发往Kafka
	logrus.Info("collect for path:%s is running...", t.path)
	for {
		line, ok := <-t.tObj.Lines //chan tail.Line
		if !ok {
			logrus.Warn("tail file close reopen,path:%s\n", t.path)
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
		msg.Topic = t.topic
		msg.Value = sarama.StringEncoder(line.Text)
		kafka.ToMsgChan(msg)
	}
}

func Init(allConf []common.CollectEntry) (err error) {
	//allConf 里面存了若干个日志的收集项
	//针对每一个日志收集项创建一个对应的tailObj

	for _, conf := range allConf {
		tt := newTailTask(conf.Path, conf.Topic) //创建一个日志收集任务
		err = tt.Init()                          //打开日志文件准备读
		if err != nil {
			logrus.Errorf("create tailObj for path:%s failed, err:%v", conf.Path, err)
			return
		}
		logrus.Info("create a tail task for path :%s success", conf.Path)
		// 起一个后台的goroutine收集日志
		go tt.run()
	}
	return
}
