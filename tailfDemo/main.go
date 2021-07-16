package main

import (
	"fmt"
	"time"

	"github.com/hpcloud/tail"
)

func main() {
	filename := `.\xx.log`
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}

	//打开文件开始读取数据
	tails, err := tail.TailFile(filename, config)
	if err != nil {
		fmt.Println("tail %s failed ,err:%v\n", filename, err)
		return
	}
	var (
		msg *tail.Line
		ok  bool
	)
	for {
		msg, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen , filename:%v\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("msg:", msg.Text)
	}
}
