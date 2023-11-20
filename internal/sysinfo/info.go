package sysinfo

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	psnet "github.com/shirou/gopsutil/v3/net"
	"net"
	"time"
)

func GetMemInfo() *mem.VirtualMemoryStat {
	mem, _ := mem.VirtualMemory()

	return mem
}

func GetCpuInfo() interface{} {
	info, _ := cpu.Info()
	percent, _ := cpu.Percent(time.Second, false)
	data := map[string]interface{}{
		"infos":   info,
		"percent": percent,
	}
	return data
}

func GetCpuLoad() *load.AvgStat {
	info, _ := load.Avg()

	return info
}

func GetHostInfo() *host.InfoStat {
	hInfo, _ := host.Info()

	return hInfo
}

func GetDiskInfo() map[string]interface{} {
	data := make(map[string]interface{}, 0)

	parts, err := disk.Partitions(true)
	if err != nil {
		fmt.Printf("get Partitions failed, err:%v\n", err)
		return data
	}
	ps := make([]map[string]interface{}, 0)
	for _, part := range parts {
		diskInfo, _ := disk.Usage(part.Mountpoint)
		ps = append(ps, map[string]interface{}{
			"part":     part.String(),
			"diskInfo": diskInfo,
		})
	}

	ioStat, _ := disk.IOCounters()
	data["infos"] = ps
	data["ioStat"] = ioStat

	return data
}

func GetNetInfo() []psnet.IOCountersStat {
	info, _ := psnet.IOCounters(true)

	return info
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String()
	}

	return ""
}
