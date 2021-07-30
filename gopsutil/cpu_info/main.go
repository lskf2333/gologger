package main

//获取cpu信息

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"time"
)

// cpu info
func getCpuInfo() {
	cpuInfos, err := cpu.Info()
	if err != nil {
		fmt.Printf("get cpu info failed, err:%v", err)
	}
	for _, ci := range cpuInfos {
		fmt.Println(ci)
	}
	// CPU使用率
	for {
		percent, _ := cpu.Percent(time.Second, false)
		fmt.Printf("cpu percent:%v\n", percent)
	}
}

//cpu 负载
func getLoad() {
	info, err := load.Avg()
	if err != nil {
		fmt.Printf("load.Avg() failed, err:%v", err)
		return
	}
	fmt.Println(info)
}

//内存信息
func getMemory() {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("mem.VirtualMemory() failed, err:%v", err)
		return
	}
	fmt.Println("mem info:%v\n", memInfo)
}

// host info
func getHostInfo() {
	hInfo, _ := host.Info()
	fmt.Printf("host info:%v uptime:%v boottime:%v\n", hInfo, hInfo.Uptime, hInfo.BootTime)
}

//磁盘信息
func getDiskInfo() {
	//获取所有的分区
	parts, err := disk.Partitions(true)
	if err != nil {
		fmt.Printf("get disk partitions failed, err:%v\n", err)
		return
	}
	fmt.Println(parts)
	for _, part := range parts {
		partInfo, err := disk.Usage(part.Mountpoint)
		if err != nil {
			fmt.Printf("get part stat failed, err:%v\n", err)
			return
		}
		fmt.Println(partInfo)
	}
}

//网络信息
func getNetInfo() {
	netIOs, err := net.IOCounters(true)
	if err != nil {
		fmt.Printf("get net io counter failed,err:%v\n", err)
		return
	}
	for _, netIO := range netIOs {
		fmt.Println(netIO)
	}
}

func main() {
	//getCpuInfo()
	//getLoad()
	//getMemory()
	//getHostInfo()
	//getDiskInfo()
	getNetInfo()
}
