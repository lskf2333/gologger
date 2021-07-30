package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	"time"
)

var (
	cli client.Client
)

func initConnInflux() (err error) {
	cli, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: "admin",
		Password: "",
	})
	return
}

func writesCpuPoints(info *CpuInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}

	//根据传入数据的类型出插入数据
	tags := map[string]string{"cpu": "cpu0"}
	fields := map[string]interface{}{
		"cpu_percent1": info.CpuPercent,
	}
	pt, err := client.NewPoint("cpu_percent1", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert cpu info success")

}

func writesMemPoints(data *MemInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}

	//根据传入数据的类型出插入数据

	tags := map[string]string{"mem": "mem"}
	fields := map[string]interface{}{
		"total":       int64(data.Total),
		"available":   int64(data.Available),
		"used":        int64(data.Used),
		"usedPercent": data.UsedPercent,
	}
	pt, err := client.NewPoint("memory", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert mem info success")
}

// cpu info
func getCpuInfo() {
	var cupInfo = new(CpuInfo)
	percent, _ := cpu.Percent(time.Second, false)
	fmt.Printf("cpu percent:%v\n", percent)
	//写入到influxDB中
	cupInfo.CpuPercent = percent[0]
	writesCpuPoints(cupInfo)
}

// 内存
func getMemory() {
	var memInfo = new(MemInfo)
	info, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("mem.VirtualMemory() failed, err:%v", err)
		return
	}
	memInfo.Total = info.Total
	memInfo.Available = info.Available
	memInfo.Used = info.Used
	memInfo.UsedPercent = info.UsedPercent
	memInfo.Buffers = info.Buffers
	memInfo.Cached = info.Cached
	writesMemPoints(memInfo)
}

func run(interval time.Duration) {
	ticker := time.Tick(interval)
	for _ = range ticker {
		getCpuInfo()
		getMemory()
	}
}

func main() {
	err := initConnInflux()
	if err != nil {
		fmt.Printf("connect to influxDB failed, err:%v\n", err)
		return
	}
	run(time.Second)
}
